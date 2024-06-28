package chainlink

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/grafana/pyroscope-go"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	pkgcapabilities "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/static"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	externalp2p "github.com/smartcontractkit/chainlink/v2/core/services/p2p/wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/promreporter"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	workflowstore "github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/ldapauth"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/localauth"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// Application implements the common functions used in the core node.
//
//go:generate mockery --quiet --name Application --output ../../internal/mocks/ --case=underscore
type Application interface {
	Start(ctx context.Context) error
	Stop() error
	GetLogger() logger.SugaredLogger
	GetAuditLogger() audit.AuditLogger
	GetHealthChecker() services.Checker
	GetDB() sqlutil.DataSource
	GetConfig() GeneralConfig
	SetLogLevel(lvl zapcore.Level) error
	GetKeyStore() keystore.Master
	WakeSessionReaper()
	GetWebAuthnConfiguration() sessions.WebAuthnConfiguration

	GetExternalInitiatorManager() webhook.ExternalInitiatorManager
	GetRelayers() RelayerChainInteroperators
	GetLoopRegistry() *plugins.LoopRegistry
	GetLoopRegistrarConfig() plugins.RegistrarConfig

	// V2 Jobs (TOML specified)
	JobSpawner() job.Spawner
	JobORM() job.ORM
	EVMORM() evmtypes.Configs
	PipelineORM() pipeline.ORM
	BridgeORM() bridges.ORM
	BasicAdminUsersORM() sessions.BasicAdminUsersORM
	AuthenticationProvider() sessions.AuthenticationProvider
	TxmStorageService() txmgr.EvmTxStore
	AddJobV2(ctx context.Context, job *job.Job) error
	DeleteJob(ctx context.Context, jobID int32) error
	RunWebhookJobV2(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta jsonserializable.JSONSerializable) (int64, error)
	ResumeJobV2(ctx context.Context, taskID uuid.UUID, result pipeline.Result) error
	// Testing only
	RunJobV2(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)

	// Feeds
	GetFeedsService() feeds.Service

	// ReplayFromBlock replays logs from on or after the given block number. If forceBroadcast is
	// set to true, consumers will reprocess data even if it has already been processed.
	ReplayFromBlock(chainID *big.Int, number uint64, forceBroadcast bool) error

	// ID is unique to this particular application instance
	ID() uuid.UUID

	SecretGenerator() SecretGenerator

	// FindLCA - finds last common ancestor for LogPoller's chain available in the database and RPC chain
	FindLCA(ctx context.Context, chainID *big.Int) (*logpoller.LogPollerBlock, error)
	// DeleteLogPollerDataAfter - delete LogPoller state starting from the specified block
	DeleteLogPollerDataAfter(ctx context.Context, chainID *big.Int, start int64) error
}

// ChainlinkApplication contains fields for the JobSubscriber, Scheduler,
// and Store. The JobSubscriber and Scheduler are also available
// in the services package, but the Store has its own package.
type ChainlinkApplication struct {
	relayers                 *CoreRelayerChainInteroperators
	jobORM                   job.ORM
	jobSpawner               job.Spawner
	pipelineORM              pipeline.ORM
	pipelineRunner           pipeline.Runner
	bridgeORM                bridges.ORM
	localAdminUsersORM       sessions.BasicAdminUsersORM
	authenticationProvider   sessions.AuthenticationProvider
	txmStorageService        txmgr.EvmTxStore
	FeedsService             feeds.Service
	webhookJobRunner         webhook.JobRunner
	Config                   GeneralConfig
	KeyStore                 keystore.Master
	ExternalInitiatorManager webhook.ExternalInitiatorManager
	SessionReaper            *utils.SleeperTask
	shutdownOnce             sync.Once
	srvcs                    []services.ServiceCtx
	HealthChecker            services.Checker
	Nurse                    *services.Nurse
	logger                   logger.SugaredLogger
	AuditLogger              audit.AuditLogger
	closeLogger              func() error
	ds                       sqlutil.DataSource
	secretGenerator          SecretGenerator
	profiler                 *pyroscope.Profiler
	loopRegistry             *plugins.LoopRegistry
	loopRegistrarConfig      plugins.RegistrarConfig

	started     bool
	startStopMu sync.Mutex
}

type ApplicationOpts struct {
	Config                     GeneralConfig
	Logger                     logger.Logger
	MailMon                    *mailbox.Monitor
	DS                         sqlutil.DataSource
	KeyStore                   keystore.Master
	RelayerChainInteroperators *CoreRelayerChainInteroperators
	AuditLogger                audit.AuditLogger
	CloseLogger                func() error
	ExternalInitiatorManager   webhook.ExternalInitiatorManager
	Version                    string
	RestrictedHTTPClient       *http.Client
	UnrestrictedHTTPClient     *http.Client
	SecretGenerator            SecretGenerator
	LoopRegistry               *plugins.LoopRegistry
	GRPCOpts                   loop.GRPCOpts
	MercuryPool                wsrpc.Pool
	CapabilitiesRegistry       coretypes.CapabilitiesRegistry
}

// NewApplication initializes a new store if one is not already
// present at the configured root directory (default: ~/.chainlink),
// the logger at the same directory and returns the Application to
// be used by the node.
// TODO: Inject more dependencies here to save booting up useless stuff in tests
func NewApplication(opts ApplicationOpts) (Application, error) {
	var srvcs []services.ServiceCtx
	auditLogger := opts.AuditLogger
	cfg := opts.Config
	relayerChainInterops := opts.RelayerChainInteroperators
	mailMon := opts.MailMon
	externalInitiatorManager := opts.ExternalInitiatorManager
	globalLogger := logger.Sugared(opts.Logger)
	keyStore := opts.KeyStore
	restrictedHTTPClient := opts.RestrictedHTTPClient
	unrestrictedHTTPClient := opts.UnrestrictedHTTPClient

	if opts.CapabilitiesRegistry == nil { // for tests only, in prod Registry is always set at this point
		opts.CapabilitiesRegistry = capabilities.NewRegistry(globalLogger)
	}

	var externalPeerWrapper p2ptypes.PeerWrapper
	var getLocalNode func(ctx context.Context) (pkgcapabilities.Node, error)
	if cfg.Capabilities().Peering().Enabled() {
		externalPeer := externalp2p.NewExternalPeerWrapper(keyStore.P2P(), cfg.Capabilities().Peering(), opts.DS, globalLogger)
		signer := externalPeer
		externalPeerWrapper = externalPeer

		srvcs = append(srvcs, externalPeerWrapper)

		dispatcher := remote.NewDispatcher(externalPeerWrapper, signer, opts.CapabilitiesRegistry, globalLogger)

		rid := cfg.Capabilities().ExternalRegistry().RelayID()
		registryAddress := cfg.Capabilities().ExternalRegistry().Address()
		relayer, err := relayerChainInterops.Get(rid)
		if err != nil {
			return nil, fmt.Errorf("could not fetch relayer %s configured for capabilities registry: %w", rid, err)
		}

		registrySyncer, err := registrysyncer.New(
			globalLogger,
			relayer,
			registryAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("could not configure syncer: %w", err)
		}

		// CEDRIC: At this point we have local node access, but we need the relayer to make calls to the CR contract
		wfLauncher := capabilities.NewLauncher(
			globalLogger,
			externalPeerWrapper,
			dispatcher,
			opts.CapabilitiesRegistry,
		)
		registrySyncer.AddLauncher(wfLauncher)

		getLocalNode = wfLauncher.LocalNode

		// We get the underlying implementation of the registry to set the GetLocalNode function.
		//
		// The rationale behind this is that the capability registry needs to have serializable methods since it's passed around to standard-capabilities as a GRPC service.
		if r, ok := opts.CapabilitiesRegistry.(*capabilities.Registry); ok {
			r.SetGetLocalNodeFunc(getLocalNode)
		} else {
			globalLogger.Warn("CapabilitiesRegistry is not a *capabilities.Registry, GetLocalNode will not be set")
		}
		srvcs = append(srvcs, dispatcher, wfLauncher, registrySyncer)
	} else {
		if r, ok := opts.CapabilitiesRegistry.(*capabilities.Registry); ok {
			// Peering isn't enabled, so we're not part of a DON.
			r.SetGetLocalNodeFunc(func(ctx context.Context) (pkgcapabilities.Node, error) {
				return pkgcapabilities.Node{}, nil
			})
		} else {
			globalLogger.Warn("CapabilitiesRegistry is not a *capabilities.Registry, GetLocalNode will not be set")
		}
	}

	// LOOPs can be created as options, in the  case of LOOP relayers, or
	// as OCR2 job implementations, in the case of Median today.
	// We will have a non-nil registry here in LOOP relayers are being used, otherwise
	// we need to initialize in case we serve OCR2 LOOPs
	loopRegistry := opts.LoopRegistry
	if loopRegistry == nil {
		loopRegistry = plugins.NewLoopRegistry(globalLogger, opts.Config.Tracing())
	}

	// If the audit logger is enabled
	if auditLogger.Ready() == nil {
		srvcs = append(srvcs, auditLogger)
	}

	var profiler *pyroscope.Profiler
	if cfg.Pyroscope().ServerAddress() != "" {
		globalLogger.Debug("Pyroscope (automatic pprof profiling) is enabled")
		var err error
		profiler, err = logger.StartPyroscope(cfg.Pyroscope(), cfg.AutoPprof())
		if err != nil {
			return nil, errors.Wrap(err, "starting pyroscope (automatic pprof profiling) failed")
		}
	} else {
		globalLogger.Debug("Pyroscope (automatic pprof profiling) is disabled")
	}

	ap := cfg.AutoPprof()
	var nurse *services.Nurse
	if ap.Enabled() {
		globalLogger.Info("Nurse service (automatic pprof profiling) is enabled")
		nurse = services.NewNurse(ap, globalLogger)
		err := nurse.Start()
		if err != nil {
			return nil, err
		}
	} else {
		globalLogger.Info("Nurse service (automatic pprof profiling) is disabled")
	}

	telemetryManager := telemetry.NewManager(cfg.TelemetryIngress(), keyStore.CSA(), globalLogger)
	srvcs = append(srvcs, telemetryManager)

	backupCfg := cfg.Database().Backup()
	if backupCfg.Mode() != config.DatabaseBackupModeNone && backupCfg.Frequency() > 0 {
		globalLogger.Infow("DatabaseBackup: periodic database backups are enabled", "frequency", backupCfg.Frequency())

		databaseBackup, err := periodicbackup.NewDatabaseBackup(cfg.Database().URL(), cfg.RootDir(), backupCfg, globalLogger)
		if err != nil {
			return nil, errors.Wrap(err, "NewApplication: failed to initialize database backup")
		}
		srvcs = append(srvcs, databaseBackup)
	} else {
		globalLogger.Info("DatabaseBackup: periodic database backups are disabled. To enable automatic backups, set Database.Backup.Mode=lite or Database.Backup.Mode=full")
	}

	// pool must be started before all relayers and stopped after them
	if opts.MercuryPool != nil {
		srvcs = append(srvcs, opts.MercuryPool)
	}

	// EVM chains are used all over the place. This will need to change for fully EVM extraction
	// TODO: BCF-2510, BCF-2511

	legacyEVMChains := relayerChainInterops.LegacyEVMChains()
	if legacyEVMChains == nil {
		return nil, fmt.Errorf("no evm chains found")
	}

	srvcs = append(srvcs, mailMon)
	srvcs = append(srvcs, relayerChainInterops.Services()...)
	promReporter := promreporter.NewPromReporter(opts.DS, legacyEVMChains, globalLogger)
	srvcs = append(srvcs, promReporter)

	// Initialize Local Users ORM and Authentication Provider specified in config
	// BasicAdminUsersORM is initialized and required regardless of separate Authentication Provider
	localAdminUsersORM := localauth.NewORM(opts.DS, cfg.WebServer().SessionTimeout().Duration(), globalLogger, auditLogger)

	// Initialize Sessions ORM based on environment configured authenticator
	// localDB auth or remote LDAP auth
	authMethod := cfg.WebServer().AuthenticationMethod()
	var authenticationProvider sessions.AuthenticationProvider
	var sessionReaper *utils.SleeperTask

	switch sessions.AuthenticationProviderName(authMethod) {
	case sessions.LDAPAuth:
		var err error
		authenticationProvider, err = ldapauth.NewLDAPAuthenticator(
			opts.DS, cfg.WebServer().LDAP(), cfg.Insecure().DevWebServer(), globalLogger, auditLogger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "NewApplication: failed to initialize LDAP Authentication module")
		}
		sessionReaper = ldapauth.NewLDAPServerStateSync(opts.DS, cfg.WebServer().LDAP(), globalLogger)
	case sessions.LocalAuth:
		authenticationProvider = localauth.NewORM(opts.DS, cfg.WebServer().SessionTimeout().Duration(), globalLogger, auditLogger)
		sessionReaper = localauth.NewSessionReaper(opts.DS, cfg.WebServer(), globalLogger)
	default:
		return nil, errors.Errorf("NewApplication: Unexpected 'AuthenticationMethod': %s supported values: %s, %s", authMethod, sessions.LocalAuth, sessions.LDAPAuth)
	}

	var (
		pipelineORM    = pipeline.NewORM(opts.DS, globalLogger, cfg.JobPipeline().MaxSuccessfulRuns())
		bridgeORM      = bridges.NewORM(opts.DS)
		mercuryORM     = mercury.NewORM(opts.DS)
		pipelineRunner = pipeline.NewRunner(pipelineORM, bridgeORM, cfg.JobPipeline(), cfg.WebServer(), legacyEVMChains, keyStore.Eth(), keyStore.VRF(), globalLogger, restrictedHTTPClient, unrestrictedHTTPClient)
		jobORM         = job.NewORM(opts.DS, pipelineORM, bridgeORM, keyStore, globalLogger)
		txmORM         = txmgr.NewTxStore(opts.DS, globalLogger)
		streamRegistry = streams.NewRegistry(globalLogger, pipelineRunner)
		workflowORM    = workflowstore.NewDBStore(opts.DS, globalLogger, clockwork.NewRealClock())
	)

	for _, chain := range legacyEVMChains.Slice() {
		chain.HeadBroadcaster().Subscribe(promReporter)
		chain.TxManager().RegisterResumeCallback(pipelineRunner.ResumeRun)
	}

	srvcs = append(srvcs, pipelineORM)

	loopRegistrarConfig := plugins.NewRegistrarConfig(opts.GRPCOpts, opts.LoopRegistry.Register, opts.LoopRegistry.Unregister)

	var (
		delegates = map[job.Type]job.Delegate{
			job.DirectRequest: directrequest.NewDelegate(
				globalLogger,
				pipelineRunner,
				pipelineORM,
				legacyEVMChains,
				mailMon),
			job.Keeper: keeper.NewDelegate(
				cfg,
				opts.DS,
				jobORM,
				pipelineRunner,
				globalLogger,
				legacyEVMChains,
				mailMon),
			job.VRF: vrf.NewDelegate(
				opts.DS,
				keyStore,
				pipelineRunner,
				pipelineORM,
				legacyEVMChains,
				globalLogger,
				mailMon),
			job.Webhook: webhook.NewDelegate(
				pipelineRunner,
				externalInitiatorManager,
				globalLogger),
			job.Cron: cron.NewDelegate(
				pipelineRunner,
				globalLogger),
			job.BlockhashStore: blockhashstore.NewDelegate(
				cfg,
				globalLogger,
				legacyEVMChains,
				keyStore.Eth()),
			job.BlockHeaderFeeder: blockheaderfeeder.NewDelegate(
				cfg,
				globalLogger,
				legacyEVMChains,
				keyStore.Eth()),
			job.Gateway: gateway.NewDelegate(
				legacyEVMChains,
				keyStore.Eth(),
				opts.DS,
				globalLogger),
			job.Stream: streams.NewDelegate(
				globalLogger,
				streamRegistry,
				pipelineRunner,
				cfg.JobPipeline(),
			),
			job.StandardCapabilities: standardcapabilities.NewDelegate(
				globalLogger,
				opts.DS, jobORM,
				opts.CapabilitiesRegistry,
				loopRegistrarConfig,
				telemetryManager,
				pipelineRunner,
				opts.RelayerChainInteroperators),
		}
		webhookJobRunner = delegates[job.Webhook].(*webhook.Delegate).WebhookJobRunner()
	)

	delegates[job.Workflow] = workflows.NewDelegate(
		globalLogger,
		opts.CapabilitiesRegistry,
		workflowORM,
		getLocalNode,
	)

	// Flux monitor requires ethereum just to boot, silence errors with a null delegate
	if !cfg.EVMRPCEnabled() {
		delegates[job.FluxMonitor] = &job.NullDelegate{Type: job.FluxMonitor}
	} else {
		delegates[job.FluxMonitor] = fluxmonitorv2.NewDelegate(
			cfg,
			keyStore.Eth(),
			jobORM,
			pipelineORM,
			pipelineRunner,
			opts.DS,
			legacyEVMChains,
			globalLogger,
		)
	}

	var peerWrapper *ocrcommon.SingletonPeerWrapper
	if !cfg.OCR().Enabled() && !cfg.OCR2().Enabled() {
		globalLogger.Debug("P2P stack not needed")
	} else if cfg.P2P().Enabled() {
		if err := ocrcommon.ValidatePeerWrapperConfig(cfg.P2P()); err != nil {
			return nil, err
		}
		peerWrapper = ocrcommon.NewSingletonPeerWrapper(keyStore, cfg.P2P(), cfg.OCR(), opts.DS, globalLogger)
		srvcs = append(srvcs, peerWrapper)
	} else {
		return nil, fmt.Errorf("P2P stack required for OCR or OCR2")
	}

	if cfg.OCR().Enabled() {
		delegates[job.OffchainReporting] = ocr.NewDelegate(
			opts.DS,
			jobORM,
			keyStore,
			pipelineRunner,
			peerWrapper,
			telemetryManager,
			legacyEVMChains,
			globalLogger,
			cfg,
			mailMon,
		)
	} else {
		globalLogger.Debug("Off-chain reporting disabled")
	}

	if cfg.OCR2().Enabled() {
		globalLogger.Debug("Off-chain reporting v2 enabled")

		ocr2DelegateConfig := ocr2.NewDelegateConfig(cfg.OCR2(), cfg.Mercury(), cfg.Threshold(), cfg.Insecure(), cfg.JobPipeline(), loopRegistrarConfig)

		delegates[job.OffchainReporting2] = ocr2.NewDelegate(
			opts.DS,
			jobORM,
			bridgeORM,
			mercuryORM,
			pipelineRunner,
			streamRegistry,
			peerWrapper,
			telemetryManager,
			legacyEVMChains,
			globalLogger,
			ocr2DelegateConfig,
			keyStore.OCR2(),
			keyStore.DKGSign(),
			keyStore.DKGEncrypt(),
			keyStore.Eth(),
			opts.RelayerChainInteroperators,
			mailMon,
			opts.CapabilitiesRegistry,
		)
		delegates[job.Bootstrap] = ocrbootstrap.NewDelegateBootstrap(
			opts.DS,
			jobORM,
			peerWrapper,
			globalLogger,
			cfg.OCR2(),
			cfg.Insecure(),
			opts.RelayerChainInteroperators,
		)
	} else {
		globalLogger.Debug("Off-chain reporting v2 disabled")
	}

	healthChecker := commonservices.NewChecker(static.Version, static.Sha)

	var lbs []utils.DependentAwaiter
	for _, c := range legacyEVMChains.Slice() {
		lbs = append(lbs, c.LogBroadcaster())
	}
	jobSpawner := job.NewSpawner(jobORM, cfg.Database(), healthChecker, delegates, globalLogger, lbs)
	srvcs = append(srvcs, jobSpawner, pipelineRunner)

	// We start the log poller after the job spawner
	// so jobs have a chance to apply their initial log filters.
	if cfg.Feature().LogPoller() {
		for _, c := range legacyEVMChains.Slice() {
			srvcs = append(srvcs, c.LogPoller())
		}
	}

	var feedsService feeds.Service
	if cfg.Feature().FeedsManager() {
		feedsORM := feeds.NewORM(opts.DS)
		feedsService = feeds.NewService(
			feedsORM,
			jobORM,
			opts.DS,
			jobSpawner,
			keyStore,
			cfg,
			cfg.Insecure(),
			cfg.JobPipeline(),
			cfg.OCR(),
			cfg.OCR2(),
			legacyEVMChains,
			globalLogger,
			opts.Version,
			loopRegistrarConfig,
		)
	} else {
		feedsService = &feeds.NullService{}
	}

	for _, s := range srvcs {
		if s == nil {
			panic("service unexpectedly nil")
		}
		if err := healthChecker.Register(s); err != nil {
			return nil, err
		}
	}

	return &ChainlinkApplication{
		relayers:                 opts.RelayerChainInteroperators,
		jobORM:                   jobORM,
		jobSpawner:               jobSpawner,
		pipelineRunner:           pipelineRunner,
		pipelineORM:              pipelineORM,
		bridgeORM:                bridgeORM,
		localAdminUsersORM:       localAdminUsersORM,
		authenticationProvider:   authenticationProvider,
		txmStorageService:        txmORM,
		FeedsService:             feedsService,
		Config:                   cfg,
		webhookJobRunner:         webhookJobRunner,
		KeyStore:                 keyStore,
		SessionReaper:            sessionReaper,
		ExternalInitiatorManager: externalInitiatorManager,
		HealthChecker:            healthChecker,
		Nurse:                    nurse,
		logger:                   globalLogger,
		AuditLogger:              auditLogger,
		closeLogger:              opts.CloseLogger,
		secretGenerator:          opts.SecretGenerator,
		profiler:                 profiler,
		loopRegistry:             loopRegistry,
		loopRegistrarConfig:      loopRegistrarConfig,

		ds: opts.DS,

		// NOTE: Can keep things clean by putting more things in srvcs instead of manually start/closing
		srvcs: srvcs,
	}, nil
}

func (app *ChainlinkApplication) SetLogLevel(lvl zapcore.Level) error {
	if err := app.Config.SetLogLevel(lvl); err != nil {
		return err
	}
	app.logger.SetLogLevel(lvl)
	return nil
}

// Start all necessary services. If successful, nil will be returned.
// Start sequence is aborted if the context gets cancelled.
func (app *ChainlinkApplication) Start(ctx context.Context) error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	if app.started {
		panic("application is already started")
	}

	if app.FeedsService != nil {
		if err := app.FeedsService.Start(ctx); err != nil {
			app.logger.Errorf("[Feeds Service] Failed to start %v", err)
			app.FeedsService = &feeds.NullService{} // so we don't try to Close() later
		}
	}

	var ms services.MultiStart
	for _, service := range app.srvcs {
		if ctx.Err() != nil {
			err := errors.Wrap(ctx.Err(), "aborting start")
			return multierr.Combine(err, ms.Close())
		}

		app.logger.Debugw("Starting service...", "name", service.Name())

		if err := ms.Start(ctx, service); err != nil {
			return err
		}
	}

	// Start HealthChecker last, so that the other services had the chance to
	// start enough to immediately pass the readiness check.
	if err := app.HealthChecker.Start(); err != nil {
		return err
	}

	app.started = true

	return nil
}

func (app *ChainlinkApplication) StopIfStarted() error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	if app.started {
		return app.stop()
	}
	return nil
}

func (app *ChainlinkApplication) GetLoopRegistry() *plugins.LoopRegistry {
	return app.loopRegistry
}
func (app *ChainlinkApplication) GetLoopRegistrarConfig() plugins.RegistrarConfig {
	return app.loopRegistrarConfig
}

// Stop allows the application to exit by halting schedules, closing
// logs, and closing the DB connection.
func (app *ChainlinkApplication) Stop() error {
	app.startStopMu.Lock()
	defer app.startStopMu.Unlock()
	return app.stop()
}

func (app *ChainlinkApplication) stop() (err error) {
	if !app.started {
		panic("application is already stopped")
	}
	app.shutdownOnce.Do(func() {
		defer func() {
			if app.closeLogger == nil {
				return
			}
			if lerr := app.closeLogger(); lerr != nil {
				err = multierr.Append(err, lerr)
			}
		}()
		app.logger.Info("Gracefully exiting...")

		// Stop services in the reverse order from which they were started
		for i := len(app.srvcs) - 1; i >= 0; i-- {
			service := app.srvcs[i]
			app.logger.Debugw("Closing service...", "name", service.Name())
			err = multierr.Append(err, service.Close())
		}

		app.logger.Debug("Stopping SessionReaper...")
		err = multierr.Append(err, app.SessionReaper.Stop())
		app.logger.Debug("Closing HealthChecker...")
		err = multierr.Append(err, app.HealthChecker.Close())
		if app.FeedsService != nil {
			app.logger.Debug("Closing Feeds Service...")
			err = multierr.Append(err, app.FeedsService.Close())
		}

		if app.Nurse != nil {
			err = multierr.Append(err, app.Nurse.Close())
		}

		if app.profiler != nil {
			err = multierr.Append(err, app.profiler.Stop())
		}

		app.logger.Info("Exited all services")

		app.started = false
	})
	return err
}

func (app *ChainlinkApplication) GetConfig() GeneralConfig {
	return app.Config
}

func (app *ChainlinkApplication) GetKeyStore() keystore.Master {
	return app.KeyStore
}

func (app *ChainlinkApplication) GetLogger() logger.SugaredLogger {
	return app.logger
}

func (app *ChainlinkApplication) GetAuditLogger() audit.AuditLogger {
	return app.AuditLogger
}

func (app *ChainlinkApplication) GetHealthChecker() services.Checker {
	return app.HealthChecker
}

func (app *ChainlinkApplication) JobSpawner() job.Spawner {
	return app.jobSpawner
}

func (app *ChainlinkApplication) JobORM() job.ORM {
	return app.jobORM
}

func (app *ChainlinkApplication) BridgeORM() bridges.ORM {
	return app.bridgeORM
}

func (app *ChainlinkApplication) BasicAdminUsersORM() sessions.BasicAdminUsersORM {
	return app.localAdminUsersORM
}

func (app *ChainlinkApplication) AuthenticationProvider() sessions.AuthenticationProvider {
	return app.authenticationProvider
}

// TODO BCF-2516 remove this all together remove EVM specifics
func (app *ChainlinkApplication) EVMORM() evmtypes.Configs {
	return app.GetRelayers().LegacyEVMChains().ChainNodeConfigs()
}

func (app *ChainlinkApplication) PipelineORM() pipeline.ORM {
	return app.pipelineORM
}

func (app *ChainlinkApplication) TxmStorageService() txmgr.EvmTxStore {
	return app.txmStorageService
}

func (app *ChainlinkApplication) GetExternalInitiatorManager() webhook.ExternalInitiatorManager {
	return app.ExternalInitiatorManager
}

func (app *ChainlinkApplication) SecretGenerator() SecretGenerator {
	return app.secretGenerator
}

// WakeSessionReaper wakes up the reaper to do its reaping.
func (app *ChainlinkApplication) WakeSessionReaper() {
	app.SessionReaper.WakeUp()
}

func (app *ChainlinkApplication) AddJobV2(ctx context.Context, j *job.Job) error {
	return app.jobSpawner.CreateJob(ctx, nil, j)
}

func (app *ChainlinkApplication) DeleteJob(ctx context.Context, jobID int32) error {
	// Do not allow the job to be deleted if it is managed by the Feeds Manager
	isManaged, err := app.FeedsService.IsJobManaged(ctx, int64(jobID))
	if err != nil {
		return err
	}

	if isManaged {
		return errors.New("job must be deleted in the feeds manager")
	}

	return app.jobSpawner.DeleteJob(ctx, nil, jobID)
}

func (app *ChainlinkApplication) RunWebhookJobV2(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta jsonserializable.JSONSerializable) (int64, error) {
	return app.webhookJobRunner.RunJob(ctx, jobUUID, requestBody, meta)
}

// Only used for local testing, not supported by the UI.
func (app *ChainlinkApplication) RunJobV2(
	ctx context.Context,
	jobID int32,
	meta map[string]interface{},
) (int64, error) {
	if build.IsProd() {
		return 0, errors.New("manual job runs not supported on secure builds")
	}
	jb, err := app.jobORM.FindJob(ctx, jobID)
	if err != nil {
		return 0, errors.Wrapf(err, "job ID %v", jobID)
	}
	var runID int64

	// Some jobs are special in that they do not have a task graph.
	isBootstrap := jb.Type == job.OffchainReporting && jb.OCROracleSpec != nil && jb.OCROracleSpec.IsBootstrapPeer
	if jb.Type.RequiresPipelineSpec() || !isBootstrap {
		var vars map[string]interface{}
		var saveTasks bool
		if jb.Type == job.VRF {
			saveTasks = true
			// Create a dummy log to trigger a run
			testLog := types.Log{
				Data: bytes.Join([][]byte{
					jb.VRFSpec.PublicKey.MustHash().Bytes(),  // key hash
					common.BigToHash(big.NewInt(42)).Bytes(), // seed
					evmutils.NewHash().Bytes(),               // sender
					evmutils.NewHash().Bytes(),               // fee
					evmutils.NewHash().Bytes()},              // requestID
					[]byte{}),
				Topics:      []common.Hash{{}, jb.ExternalIDEncodeBytesToTopic()}, // jobID BYTES
				TxHash:      evmutils.NewHash(),
				BlockNumber: 10,
				BlockHash:   evmutils.NewHash(),
			}
			vars = map[string]interface{}{
				"jobSpec": map[string]interface{}{
					"databaseID":    jb.ID,
					"externalJobID": jb.ExternalJobID,
					"name":          jb.Name.ValueOrZero(),
					"publicKey":     jb.VRFSpec.PublicKey[:],
					"evmChainID":    jb.VRFSpec.EVMChainID.String(),
				},
				"jobRun": map[string]interface{}{
					"meta":           meta,
					"logBlockHash":   testLog.BlockHash[:],
					"logBlockNumber": testLog.BlockNumber,
					"logTxHash":      testLog.TxHash,
					"logTopics":      testLog.Topics,
					"logData":        testLog.Data,
				},
			}
		} else {
			vars = map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": meta,
				},
			}
		}
		runID, _, err = app.pipelineRunner.ExecuteAndInsertFinishedRun(ctx, *jb.PipelineSpec, pipeline.NewVarsFrom(vars), app.logger, saveTasks)
	}
	return runID, err
}

func (app *ChainlinkApplication) ResumeJobV2(
	ctx context.Context,
	taskID uuid.UUID,
	result pipeline.Result,
) error {
	return app.pipelineRunner.ResumeRun(ctx, taskID, result.Value, result.Error)
}

func (app *ChainlinkApplication) GetFeedsService() feeds.Service {
	return app.FeedsService
}

// ReplayFromBlock implements the Application interface.
func (app *ChainlinkApplication) ReplayFromBlock(chainID *big.Int, number uint64, forceBroadcast bool) error {
	chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		return err
	}
	chain.LogBroadcaster().ReplayFromBlock(int64(number), forceBroadcast)
	if app.Config.Feature().LogPoller() {
		chain.LogPoller().ReplayAsync(int64(number))
	}
	return nil
}

func (app *ChainlinkApplication) GetRelayers() RelayerChainInteroperators {
	return app.relayers
}

func (app *ChainlinkApplication) GetDB() sqlutil.DataSource {
	return app.ds
}

// Returns the configuration to use for creating and authenticating
// new WebAuthn credentials
func (app *ChainlinkApplication) GetWebAuthnConfiguration() sessions.WebAuthnConfiguration {
	rpid := app.Config.WebServer().MFA().RPID()
	rporigin := app.Config.WebServer().MFA().RPOrigin()
	if rpid == "" {
		app.GetLogger().Errorf("RPID is not set, WebAuthn will likely not work as intended")
	}

	if rporigin == "" {
		app.GetLogger().Errorf("RPOrigin is not set, WebAuthn will likely not work as intended")
	}

	return sessions.WebAuthnConfiguration{
		RPID:     rpid,
		RPOrigin: rporigin,
	}
}

func (app *ChainlinkApplication) ID() uuid.UUID {
	return app.Config.AppID()
}

// FindLCA - finds last common ancestor
func (app *ChainlinkApplication) FindLCA(ctx context.Context, chainID *big.Int) (*logpoller.LogPollerBlock, error) {
	chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		return nil, err
	}
	if !app.Config.Feature().LogPoller() {
		return nil, fmt.Errorf("FindLCA is only available if LogPoller is enabled")
	}

	lca, err := chain.LogPoller().FindLCA(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find lca: %w", err)
	}

	return lca, nil
}

// DeleteLogPollerDataAfter - delete LogPoller state starting from the specified block
func (app *ChainlinkApplication) DeleteLogPollerDataAfter(ctx context.Context, chainID *big.Int, start int64) error {
	chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		return err
	}
	if !app.Config.Feature().LogPoller() {
		return fmt.Errorf("DeleteLogPollerDataAfter is only available if LogPoller is enabled")
	}

	err = chain.LogPoller().DeleteLogsAndBlocksAfter(ctx, start)
	if err != nil {
		return fmt.Errorf("failed to recover LogPoller: %w", err)
	}

	return nil
}
