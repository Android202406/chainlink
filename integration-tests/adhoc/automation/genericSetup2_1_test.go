package automation

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/stretchr/testify/require"

	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/automationv2"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	contractseth "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_triggered_streams_lookup_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
)

var (
	baseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
	secretsTOML = `[Mercury.Credentials.%s]
LegacyURL = '%s'
URL = '%s'
Username = '%s'
Password = '%s'`

	nodeSpec = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "4Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "4Gi",
			},
		},
	}

	dbSpec = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "4000m",
				"memory": "4Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "4000m",
				"memory": "4Gi",
			},
		},
		"stateful": true,
		"capacity": "10Gi",
	}
)

func setupEnvironment(t *testing.T, loadedTestConfig tc.TestConfig) (*automationv2.AutomationTest, *environment.Environment) {
	l := logging.GetTestLogger(t)

	testNetwork := networks.MustGetSelectedNetworkConfig(loadedTestConfig.Network)[0]

	testEnvironment := environment.New(&environment.Config{
		TTL: time.Hour * 24 * 30,
		NamespacePrefix: fmt.Sprintf(
			"automation-%s-%s",
			"adhoc",
			strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"),
		),
		Test:               t,
		PreventPodEviction: true,
	})

	testEnvironment.
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
			Values: map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "8000m",
						"memory": "8Gi",
					},
					"limits": map[string]interface{}{
						"cpu":    "16000m",
						"memory": "16Gi",
					},
				},
				"geth": map[string]interface{}{
					"blocktime": *loadedTestConfig.Automation.General.BlockTime,
					"capacity":  "20Gi",
				},
			},
		}))

	if testNetwork.Simulated {
		err := testEnvironment.Run()
		require.NoError(t, err, "Error launching test environment")
	}

	if testEnvironment.WillUseRemoteRunner() {
		return nil, testEnvironment
	}

	if *loadedTestConfig.Pyroscope.Enabled {
		loadedTestConfig.Pyroscope.Environment = &testEnvironment.Cfg.Namespace
	}

	if *loadedTestConfig.Automation.Adhoc.ConnectDataStream {
		secretsTOML = fmt.Sprintf(
			secretsTOML, "cred1",
			*loadedTestConfig.Automation.Adhoc.DataStreamURL, *loadedTestConfig.Automation.Adhoc.DataStreamURL,
			*loadedTestConfig.Automation.Adhoc.DataStreamUsername, *loadedTestConfig.Automation.Adhoc.DataStreamPassword,
		)
	} else {
		secretsTOML = ""
	}

	numberOfNodes := *loadedTestConfig.Automation.General.NumberOfNodes
	l.Info().Int("Number of Nodes", numberOfNodes).Msg("Number of Nodes")

	for i := 0; i < numberOfNodes+1; i++ { // +1 for the OCR boot node
		var nodeTOML string
		nodeTOML = fmt.Sprintf("%s\n\n[Log]\nLevel = \"%s\"", baseTOML, *loadedTestConfig.Automation.General.ChainlinkNodeLogLevel)
		nodeTOML = networks.AddNetworksConfig(nodeTOML, loadedTestConfig.Pyroscope, testNetwork)

		var overrideFn = func(_ interface{}, target interface{}) {
			ctfconfig.MustConfigOverrideChainlinkVersion(loadedTestConfig.ChainlinkImage, target)
			ctfconfig.MightConfigOverridePyroscopeKey(loadedTestConfig.Pyroscope, target)
		}

		cd := chainlink.NewWithOverride(i, map[string]any{
			"toml":        nodeTOML,
			"chainlink":   nodeSpec,
			"db":          dbSpec,
			"prometheus":  *loadedTestConfig.Automation.General.UsePrometheus,
			"secretsToml": secretsTOML,
		}, loadedTestConfig.ChainlinkImage, overrideFn)

		testEnvironment.AddHelm(cd)
	}

	err := testEnvironment.Run()
	require.NoError(t, err, "Error running chainlink DON")

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment, l)
	require.NoError(t, err, "Error building chain client")

	contractDeployer, err := contracts.NewContractDeployer(chainClient, l)
	require.NoError(t, err, "Error building contract deployer")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to chainlink nodes")

	chainClient.ParallelTransactions(true)

	automationTest := automationv2.NewAutomationTestK8s(chainClient, contractDeployer, chainlinkNodes)

	return automationTest, testEnvironment
}

func TestAutomation(t *testing.T) {
	l := logging.GetTestLogger(t)
	loadedTestConfig, err := tc.GetConfig("Adhoc", tc.Automation)
	if err != nil {
		t.Fatal(err)
	}
	l.Info().Interface("loadedTestConfig", loadedTestConfig).Msg("Loaded test config")

	automationTest, testEnvironment := setupEnvironment(t, loadedTestConfig)

	automationDefaultLinkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(10000))) //10000 LINK
	automationTest.RegistrySettings = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(0),
		FlatFeeMicroLINK:     uint32(40_000),
		BlockCountPerTurn:    big.NewInt(100),
		CheckGasLimit:        uint32(50_000_000), //45M
		StalenessSeconds:     big.NewInt(90_000),
		GasCeilingMultiplier: uint16(2),
		MaxPerformGas:        uint32(6_000_000),
		MinUpkeepSpend:       big.NewInt(0),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
		MaxCheckDataSize:     uint32(5_000),
		MaxPerformDataSize:   uint32(5_000),
		RegistryVersion:      contractseth.RegistryVersion_2_1,
	}
	automationTest.RegistrarSettings = contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: uint8(2),
		AutoApproveMaxAllowed: math.MaxUint16,
		MinLinkJuels:          big.NewInt(0),
	}
	automationTest.PluginConfig = ocr2keepers30config.OffchainConfig{
		TargetProbability:    "0.999",
		TargetInRounds:       1,
		PerformLockoutWindow: 75_000,
		GasLimitPerReport:    5_300_000,
		GasOverheadPerUpkeep: 300_000,
		MinConfirmations:     0,
		MaxUpkeepBatchSize:   10,
	}
	automationTest.PublicConfig = ocr3.PublicConfig{
		DeltaProgress:                           10 * time.Second,
		DeltaResend:                             15 * time.Second,
		DeltaInitial:                            500 * time.Millisecond,
		DeltaRound:                              1000 * time.Millisecond,
		DeltaGrace:                              200 * time.Millisecond,
		DeltaCertifiedCommitRequest:             300 * time.Millisecond,
		DeltaStage:                              25 * time.Second,
		RMax:                                    24,
		MaxDurationQuery:                        20 * time.Millisecond,
		MaxDurationObservation:                  20 * time.Millisecond,
		MaxDurationShouldAcceptAttestedReport:   1200 * time.Millisecond,
		MaxDurationShouldTransmitAcceptedReport: 20 * time.Millisecond,
		F:                                       1,
	}

	config := loadedTestConfig.Automation.Adhoc

	if *config.ConnectDataStream {
		automationTest.MercuryCredentialName = "cred1"
	}

	if *config.DeployContracts {
		automationTest.SetupAutomationDeployment(t)
		err = actions.FundChainlinkNodesAddress(
			automationTest.ChainlinkNodesk8s[1:], automationTest.ChainClient,
			big.NewFloat(*loadedTestConfig.Common.ChainlinkNodeFunding), 0)
		require.NoError(t, err, "Error funding chainlink nodes")
	}

	if *config.LoadContracts {
		if *config.DeleteExistingJobs {
			err = automationTest.CollectNodeDetails()
			require.NoError(t, err, "Error collecting node details")
			err = actions.DeleteAllJobs(automationTest.ChainlinkNodesk8s)
			require.NoError(t, err, "Error deleting all jobs")
		}
		automationTest.LoadAutomationDeployment(t, *config.LinkTokenAddress, *config.NativeLinkFeedAddress,
			*config.FastGasFeedAddress, *config.TranscoderAddress, *config.RegistryAddress, *config.RegistrarAddress)
	}

	if *config.SetupUpkeeps {
		upkeepConfigs := make([]automationv2.UpkeepConfig, 0)

		var bytes0 = [32]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}

		utilsABI, err := automation_utils_2_1.AutomationUtilsMetaData.GetAbi()
		require.NoError(t, err, "Error getting automation utils abi")
		logTriggerUpkeepABI, err := log_upkeep_counter_wrapper.LogUpkeepCounterMetaData.GetAbi()
		require.NoError(t, err, "Error getting log upkeep abi")
		logTriggerStreamsLookupUpkeepABI, err := log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookupMetaData.GetAbi()
		require.NoError(t, err, "Error getting log triggered streams lookup upkeep abi")

		var conditionalConsumerAddress string

		if *config.ConnectDataStream {
			conditionalConsumer, err := automationTest.Deployer.DeployAutomationStreamsLookupUpkeepConsumer(big.NewInt(math.MaxInt64), big.NewInt(10), false, true, false)
			require.NoError(t, err, "Error deploying streams lookup conditional consumer")
			conditionalConsumerAddress = conditionalConsumer.Address()

			err = conditionalConsumer.SetFeeds([]string{*loadedTestConfig.Automation.Adhoc.DataStreamFeedId})
			require.NoError(t, err, "Error setting feeds")
		} else {
			conditionalConsumer, err := automationTest.Deployer.DeployUpkeepCounter(big.NewInt(math.MaxInt64), big.NewInt(10))
			require.NoError(t, err, "Error deploying conditional consumer")
			conditionalConsumerAddress = conditionalConsumer.Address()
		}

		upkeepConfig := automationv2.UpkeepConfig{
			UpkeepName:     fmt.Sprintf("ConditionalUpkeep-%d", 0),
			EncryptedEmail: []byte("test@mail.com"),
			UpkeepContract: common.HexToAddress(conditionalConsumerAddress),
			GasLimit:       1_000_000,
			AdminAddress:   common.HexToAddress(automationTest.ChainClient.GetDefaultWallet().Address()),
			TriggerType:    uint8(0),
			CheckData:      []byte("0"),
			TriggerConfig:  []byte("0"),
			OffchainConfig: []byte("0"),
			FundingAmount:  automationDefaultLinkFunds,
		}
		l.Debug().Interface("Upkeep Config", upkeepConfig).Msg("Conditional Upkeep Config")
		upkeepConfigs = append(upkeepConfigs, upkeepConfig)

		logTriggerConsumer := *new(contracts.KeeperConsumer)
		var topic0Match [32]byte
		if *config.ConnectDataStream {
			logTriggerConsumerA, err := automationTest.Deployer.DeployAutomationLogTriggeredStreamsLookupUpkeepConsumer()
			require.NoError(t, err, "Error deploying log trigger consumer")
			err = logTriggerConsumerA.SetFeeds([]string{*loadedTestConfig.Automation.Adhoc.DataStreamFeedId})
			require.NoError(t, err, "Error setting feeds")
			logTriggerConsumer = logTriggerConsumerA
			topic0Match = logTriggerStreamsLookupUpkeepABI.Events["LimitOrderExecuted"].ID
		} else {
			logTriggerConsumer, err = automationTest.Deployer.DeployAutomationLogTriggerConsumer(big.NewInt(1))
			require.NoError(t, err, "Error deploying log trigger consumer")
			topic0Match = logTriggerUpkeepABI.Events["Trigger"].ID
		}

		logTriggerConfigStruct := automation_utils_2_1.LogTriggerConfig{
			ContractAddress: common.HexToAddress(logTriggerConsumer.Address()),
			FilterSelector:  0,
			Topic0:          topic0Match,
			Topic1:          bytes0,
			Topic2:          bytes0,
			Topic3:          bytes0,
		}
		encodedLogTriggerConfig, err := utilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
		require.NoError(t, err, "Error encoding log trigger config")
		l.Debug().Bytes("Encoded Log Trigger Config", encodedLogTriggerConfig).Msg("Encoded Log Trigger Config")

		upkeepConfig = automationv2.UpkeepConfig{
			UpkeepName:     fmt.Sprintf("LogTriggerUpkeep-%d", 0),
			EncryptedEmail: []byte("test@mail.com"),
			UpkeepContract: common.HexToAddress(logTriggerConsumer.Address()),
			GasLimit:       1_000_000,
			AdminAddress:   common.HexToAddress(automationTest.ChainClient.GetDefaultWallet().Address()),
			TriggerType:    uint8(1),
			CheckData:      []byte("0"),
			TriggerConfig:  encodedLogTriggerConfig,
			OffchainConfig: []byte("0"),
			FundingAmount:  automationDefaultLinkFunds,
		}
		l.Debug().Interface("Upkeep Config", upkeepConfig).Msg("LogTrigger Upkeep Config")
		upkeepConfigs = append(upkeepConfigs, upkeepConfig)

		registrationTxHashes, err := automationTest.RegisterUpkeeps(upkeepConfigs)
		require.NoError(t, err, "Error registering upkeeps")

		err = automationTest.ChainClient.WaitForEvents()
		require.NoError(t, err, "Failed waiting for upkeeps to register")

		upkeepIds, err := automationTest.ConfirmUpkeepsRegistered(registrationTxHashes)
		require.NoError(t, err, "Error confirming upkeeps registered")

		l.Info().Msg("Successfully registered all Automation Upkeeps")
		l.Info().Interface("Upkeep IDs", upkeepIds).Msg("Upkeeps Registered")

		l.Info().Msg("Waiting 30s for config to be picked/ up")
		time.Sleep(30 * time.Second)
		err = logTriggerConsumer.Start()
		require.NoError(t, err, "Error starting log trigger consumer")
	}

	if *config.TearDownDeployment {
		err = automationTest.CollectNodeDetails()
		require.NoError(t, err, "Error collecting node details")
		err = actions.DeleteAllJobs(automationTest.ChainlinkNodesk8s)
		require.NoError(t, err, "Error deleting all jobs")
		err = actions.ReturnFunds(automationTest.ChainlinkNodesk8s, automationTest.ChainClient)
		require.NoError(t, err, "Error returning funds")
		err = testEnvironment.Shutdown()
		require.NoError(t, err, "Error shutting down test environment")
	}

}
