package evm

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	applicationJson     = "application/json"
	blockNumber         = "blockNumber" // valid for v0.2
	feedIDs             = "feedIDs"     // valid for v0.3
	feedIdHex           = "feedIdHex"   // valid for v0.2
	headerAuthorization = "Authorization"
	headerContentType   = "Content-Type"
	headerTimestamp     = "X-Authorization-Timestamp"
	headerSignature     = "X-Authorization-Signature-SHA256"
	headerUpkeepId      = "X-Authorization-Upkeep-Id"
	mercuryPathV02      = "/client?"              // only used to access mercury v0.2 server
	mercuryBatchPathV03 = "/api/v1/reports/bulk?" // only used to access mercury v0.3 server
	retryDelay          = 500 * time.Millisecond
	timestamp           = "timestamp" // valid for v0.3
	totalAttempt        = 3
)

type StreamsLookup struct {
	feedParamKey string
	feeds        []string
	timeParamKey string
	time         *big.Int
	extraData    []byte
	upkeepId     *big.Int
	block        uint64
}

type ComposerRequestV1 struct {
	scriptHash         string
	functionsArguments []string
	useMercury         bool
	feedParamKey       string
	feeds              []string
	timeParamKey       string
	time               *big.Int
	extraData          []byte
	upkeepId           *big.Int
	block              uint64
}

type MercuryReportStruct struct {
	ObservationsTimestamp uint32   `json:"observations_timestamp"`
	Price                 *big.Int `json:"price"`
	Bid                   *big.Int `json:"bid"`
	Ask                   *big.Int `json:"ask"`
}

// MercuryV02Response represents a JSON structure used by Mercury v0.2
type MercuryV02Response struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

// MercuryV03Response represents a JSON structure used by Mercury v0.3
type MercuryV03Response struct {
	Reports []MercuryV03Report `json:"reports"`
}

type MercuryV03Report struct {
	FeedID                []byte `json:"feedID"` // feed id in hex
	ValidFromTimestamp    uint32 `json:"validFromTimestamp"`
	ObservationsTimestamp uint32 `json:"observationsTimestamp"`
	FullReport            []byte `json:"fullReport"` // the actual mercury report of this feed, can be sent to verifier
}

type MercuryData struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     [][]byte
	State     encoding.PipelineExecutionState
}

// UpkeepPrivilegeConfig represents the administrative offchain config for each upkeep. It can be set by s_upkeepPrivilegeManager
// role on the registry. Upkeeps allowed to use Mercury server will have this set to true.
type UpkeepPrivilegeConfig struct {
	MercuryEnabled bool `json:"mercuryEnabled"`
}

// streamsLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) streamsLookup(ctx context.Context, checkResults []ocr2keepers.CheckResult) []ocr2keepers.CheckResult {
	lggr := r.lggr.With("where", "StreamsLookup")
	lookups := map[int]*StreamsLookup{}
	for i, res := range checkResults {
		if res.IneligibilityReason != uint8(encoding.UpkeepFailureReasonTargetCheckReverted) {
			// Streams Lookup only works when upkeep target check reverts
			continue
		}

		block := big.NewInt(int64(res.Trigger.BlockNumber))
		upkeepId := res.UpkeepID

		// Try to decode the revert error into streams lookup format. User upkeeps can revert with any reason, see if they
		// tried to call mercury
		lggr.Infof("at block %d upkeep %s trying to decodeStreamsLookup performData=%s", block, upkeepId, hexutil.Encode(checkResults[i].PerformData))
		l, err := r.decodeStreamsLookup(res.PerformData)
		if err != nil {
			lggr.Warnf("at block %d upkeep %s decodeStreamsLookup failed: %v", block, upkeepId, err)
			// user contract did not revert with StreamsLookup error
			continue
		}
		if r.mercury.cred == nil {
			lggr.Errorf("at block %d upkeep %s tries to access mercury server but mercury credential is not configured", block, upkeepId)
			continue
		}

		if len(l.feeds) == 0 {
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			lggr.Warnf("at block %s upkeep %s has empty feeds array", block, upkeepId)
			continue
		}
		// mercury permission checking for v0.3 is done by mercury server
		if l.feedParamKey == feedIdHex && l.timeParamKey == blockNumber {
			// check permission on the registry for mercury v0.2
			opts := r.buildCallOpts(ctx, block)
			state, reason, retryable, allowed, err := r.allowedToUseMercury(opts, upkeepId.BigInt())
			if err != nil {
				lggr.Warnf("at block %s upkeep %s failed to query mercury allow list: %s", block, upkeepId, err)
				checkResults[i].PipelineExecutionState = uint8(state)
				checkResults[i].IneligibilityReason = uint8(reason)
				checkResults[i].Retryable = retryable
				continue
			}

			if !allowed {
				lggr.Warnf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
				checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryAccessNotAllowed)
				continue
			}
		} else if l.feedParamKey != feedIDs || l.timeParamKey != timestamp {
			// if mercury version cannot be determined, set failure reason
			lggr.Warnf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			continue
		}

		l.upkeepId = upkeepId.BigInt()
		// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
		// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
		l.block = uint64(block.Int64())
		lggr.Infof("at block %d upkeep %s decodeStreamsLookup feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", block, upkeepId, l.feedParamKey, l.timeParamKey, l.feeds, l.time, hexutil.Encode(l.extraData))
		lookups[i] = l
	}

	var wg sync.WaitGroup
	for i, lookup := range lookups {
		wg.Add(1)
		go r.doLookup(ctx, &wg, lookup, i, checkResults, lggr)
	}
	wg.Wait()

	// don't surface error to plugin bc StreamsLookup process should be self-contained.
	return checkResults
}

func (r *EvmRegistry) doLookup(ctx context.Context, wg *sync.WaitGroup, lookup *StreamsLookup, i int, checkResults []ocr2keepers.CheckResult, lggr logger.Logger) {
	defer wg.Done()

	state, reason, values, retryable, err := r.doMercuryRequest(ctx, lookup, lggr)
	if err != nil {
		lggr.Errorf("upkeep %s retryable %v doMercuryRequest: %s", lookup.upkeepId, retryable, err.Error())
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		checkResults[i].IneligibilityReason = uint8(reason)
		return
	}

	for j, v := range values {
		lggr.Infof("upkeep %s doMercuryRequest values[%d]: %s", lookup.upkeepId, j, hexutil.Encode(v))
	}

	state, retryable, mercuryBytes, err := r.checkCallback(ctx, values, lookup)
	if err != nil {
		lggr.Errorf("at block %d upkeep %s checkCallback err: %s", lookup.block, lookup.upkeepId, err.Error())
		checkResults[i].Retryable = retryable
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}
	lggr.Infof("checkCallback mercuryBytes=%s", hexutil.Encode(mercuryBytes))

	state, needed, performData, failureReason, _, err := r.packer.UnpackCheckCallbackResult(mercuryBytes)
	if err != nil {
		lggr.Errorf("at block %d upkeep %s UnpackCheckCallbackResult err: %s", lookup.block, lookup.upkeepId, err.Error())
		checkResults[i].PipelineExecutionState = uint8(state)
		return
	}

	if failureReason == uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted) {
		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryCallbackReverted)
		lggr.Debugf("at block %d upkeep %s mercury callback reverts", lookup.block, lookup.upkeepId)
		return
	}

	if !needed {
		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonUpkeepNotNeeded)
		lggr.Debugf("at block %d upkeep %s callback reports upkeep not needed", lookup.block, lookup.upkeepId)
		return
	}

	checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonNone)
	checkResults[i].Eligible = true
	checkResults[i].PerformData = performData
	lggr.Infof("at block %d upkeep %s successful with perform data: %s", lookup.block, lookup.upkeepId, hexutil.Encode(performData))
}

// allowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
func (r *EvmRegistry) allowedToUseMercury(opts *bind.CallOpts, upkeepId *big.Int) (state encoding.PipelineExecutionState, reason encoding.UpkeepFailureReason, retryable bool, allow bool, err error) {
	allowed, ok := r.mercury.allowListCache.Get(upkeepId.String())
	if ok {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonNone, false, allowed.(bool), nil
	}

	payload, err := r.packer.PackGetUpkeepPrivilegeConfig(upkeepId)
	if err != nil {
		// pack error, no retryable
		r.lggr.Warnf("failed to pack getUpkeepPrivilegeConfig data for upkeepId %s: %s", upkeepId, err)

		return encoding.PackUnpackDecodeFailed, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to pack upkeepId: %w", err)
	}

	var resultBytes hexutil.Bytes
	args := map[string]interface{}{
		"to":   r.addr.Hex(),
		"data": hexutil.Bytes(payload),
	}

	// call checkCallback function at the block which OCR3 has agreed upon
	err = r.client.CallContext(opts.Context, &resultBytes, "eth_call", args, hexutil.EncodeUint64(opts.BlockNumber.Uint64()))
	if err != nil {
		return encoding.RpcFlakyFailure, encoding.UpkeepFailureReasonNone, true, false, fmt.Errorf("failed to get upkeep privilege config: %v", err)
	}

	cfg, err := r.packer.UnpackGetUpkeepPrivilegeConfig(resultBytes)
	if err != nil {
		return encoding.PackUnpackDecodeFailed, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to get upkeep privilege config: %v", err)
	}

	if len(cfg) == 0 {
		r.mercury.allowListCache.Set(upkeepId.String(), false, cache.DefaultExpiration)
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonMercuryAccessNotAllowed, false, false, fmt.Errorf("upkeep privilege config is empty")
	}

	var privilegeConfig UpkeepPrivilegeConfig
	if err := json.Unmarshal(cfg, &privilegeConfig); err != nil {
		return encoding.MercuryUnmarshalError, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to unmarshal privilege config: %v", err)
	}

	r.mercury.allowListCache.Set(upkeepId.String(), privilegeConfig.MercuryEnabled, cache.DefaultExpiration)

	return encoding.NoPipelineError, encoding.UpkeepFailureReasonNone, false, privilegeConfig.MercuryEnabled, nil
}

// decodeStreamsLookup decodes the revert error StreamsLookup(string feedParamKey, string[] feeds, string feedParamKey, uint256 time, byte[] extraData)
func (r *EvmRegistry) decodeStreamsLookup(data []byte) (*StreamsLookup, error) {
	e := r.mercury.abi.Errors["StreamsLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	return &StreamsLookup{
		feedParamKey: *abi.ConvertType(errorParameters[0], new(string)).(*string),
		feeds:        *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		timeParamKey: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		time:         *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		extraData:    *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}, nil
}

// streamsLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) composerRequest(ctx context.Context, checkResults []ocr2keepers.CheckResult) []ocr2keepers.CheckResult {
	lggr := r.lggr.With("where", "ComposerRequest")
	requests := map[int]*ComposerRequestV1{}
	for i, res := range checkResults {
		if res.IneligibilityReason != uint8(encoding.UpkeepFailureReasonTargetCheckReverted) {
			// Streams Lookup only works when upkeep target check reverts
			continue
		}

		block := big.NewInt(int64(res.Trigger.BlockNumber))
		upkeepId := res.UpkeepID

		// Try to decode the revert error into streams lookup format. User upkeeps can revert with any reason, see if they
		// tried to call mercury
		lggr.Infof("at block %d upkeep %s trying to decodeComposerRequest performData=%s", block, upkeepId, hexutil.Encode(checkResults[i].PerformData))
		req, err := r.decodeComposerRequest(res.PerformData)
		if err != nil {
			lggr.Warnf("at block %d upkeep %s decodeComposerRequest failed: %v", block, upkeepId, err)
			// user contract did not revert with StreamsLookup error
			continue
		}
		if r.mercury.cred == nil && req.useMercury {
			lggr.Errorf("at block %d upkeep %s tries to access mercury server for composer request but mercury credential is not configured", block, upkeepId)
			continue
		}

		if len(req.feeds) == 0 && req.useMercury {
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			lggr.Warnf("at block %s upkeep %s has empty feeds array", block, upkeepId)
			continue
		}
		// mercury permission checking for v0.3 is done by mercury server
		if req.feedParamKey == feedIdHex && req.timeParamKey == blockNumber && req.useMercury {
			// check permission on the registry for mercury v0.2
			opts := r.buildCallOpts(ctx, block)
			state, reason, retryable, allowed, err := r.allowedToUseMercury(opts, upkeepId.BigInt())
			if err != nil {
				lggr.Warnf("at block %s upkeep %s failed to query mercury allow list: %s", block, upkeepId, err)
				checkResults[i].PipelineExecutionState = uint8(state)
				checkResults[i].IneligibilityReason = uint8(reason)
				checkResults[i].Retryable = retryable
				continue
			}

			if !allowed {
				lggr.Warnf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
				checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryAccessNotAllowed)
				continue
			}
		} else if (req.feedParamKey != feedIDs || req.timeParamKey != timestamp) && req.useMercury {
			// if mercury version cannot be determined, set failure reason
			lggr.Warnf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			continue
		}

		req.upkeepId = upkeepId.BigInt()
		// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
		// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
		req.block = uint64(block.Int64())
		lggr.Infof("at block %d upkeep %s decodeStreamsLookup feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", block, upkeepId, req.feedParamKey, req.timeParamKey, req.feeds, req.time, hexutil.Encode(req.extraData))
		requests[i] = req
	}

	var wg sync.WaitGroup
	for i, req := range requests {
		wg.Add(1)
		go r.doComposerRequest(ctx, &wg, req, i, checkResults, lggr)
	}
	wg.Wait()

	// don't surface error to plugin bc StreamsLookup process should be self-contained.
	return checkResults
}

func (r *EvmRegistry) doComposerRequest(ctx context.Context, wg *sync.WaitGroup, request *ComposerRequestV1, i int, checkResults []ocr2keepers.CheckResult, lggr logger.Logger) {
	defer wg.Done()

	var state encoding.PipelineExecutionState
	var reason encoding.UpkeepFailureReason
	var values [][]byte
	var retryable bool
	var err error

	if request.useMercury {
		lookup := &StreamsLookup{
			feedParamKey: request.feedParamKey,
			feeds:        request.feeds,
			timeParamKey: request.timeParamKey,
			time:         request.time,
			extraData:    request.extraData,
			upkeepId:     request.upkeepId,
			block:        request.block,
		}
		state, reason, values, retryable, err = r.doMercuryRequest(ctx, lookup, lggr)
		if err != nil {
			lggr.Errorf("upkeep %s retryable %v doMercuryRequest: %s", request.upkeepId, retryable, err.Error())
			checkResults[i].Retryable = retryable
			checkResults[i].PipelineExecutionState = uint8(state)
			checkResults[i].IneligibilityReason = uint8(reason)
			return
		}

		var reportData []MercuryReportStruct
		for j, v := range values {
			lggr.Infof("upkeep %s doMercuryRequest values[%d]: %s", request.upkeepId, j, hexutil.Encode(v))
			const innerReportType = `[{"type":"bytes32"},{"type":"uint32"},{"type":"int192"},{"type":"int192"},{"type":"int192"},{"type":"uint64"},{"type":"bytes32"},{"type":"uint64"},{"type":"uint64"}]`
			const reportType = `[{ "type": "bytes32[3]" },{ "type": "bytes" },{ "type": "bytes32[]" },{ "type": "bytes32[]" },{ "type": "bytes32" }]`
			interfaces, err := utils.ABIDecode(reportType, v)
			if err != nil {
				panic(err)
			}
			innerInterfaces, err := utils.ABIDecode(innerReportType, interfaces[1].([]byte))
			if err != nil {
				panic(err)
			}
			fmt.Println("DATA_FROM_VALUE")
			fmt.Println(common.Hash(innerInterfaces[0].([32]byte)).String())
			fmt.Println(innerInterfaces[1].(uint32))
			fmt.Println(innerInterfaces[2].(*big.Int))
			fmt.Println(innerInterfaces[3].(*big.Int))
			fmt.Println(innerInterfaces[4].(*big.Int))
			fmt.Println(innerInterfaces[5].(uint64))
			reportData = append(reportData, MercuryReportStruct{
				ObservationsTimestamp: innerInterfaces[1].(uint32),
				Price:                 innerInterfaces[2].(*big.Int),
				Bid:                   innerInterfaces[3].(*big.Int),
				Ask:                   innerInterfaces[4].(*big.Int),
			})
		}

		mercuryArg, err := json.Marshal(reportData)
		if err != nil {
			fmt.Println("ERROR MARSHALLING DATA", err)
		}

		// THIS PRINTS OUR MERCURY ARGUMENT AND THE ON-CHAIN PRICE DATA.
		fmt.Println("MERCURY_ARG", string(mercuryArg), request.functionsArguments)

		// SIM START
		// AT THIS POINT, FUNCTIONS SHOULD BE CALLED VIA GATEWAY TO DO SOMETHING WITH THE REVERT ARGUMENTS AND MERCURY DATA.
		// UNTIL THEN, THE MERCURY REPORTS ARE SIMPLY PASSED BACK AS PERFORM DATA.

		abiEncode, err := utils.ABIEncode(`[{"type":"bytes[]"},{"type":"bytes"}]`, values, lookup.extraData)
		if err != nil {
			lggr.Errorf("upkeep %s retryable %v abi encode packing: %s", request.upkeepId, retryable, err.Error())
			checkResults[i].Retryable = false
			checkResults[i].PipelineExecutionState = uint8(encoding.PackUnpackDecodeFailed)
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
			return
		}

		// SIM END

		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonNone)
		checkResults[i].Eligible = true
		checkResults[i].PerformData = abiEncode
	}
}

// decodeStreamsLookup decodes the revert error StreamsLookup(string feedParamKey, string[] feeds, string feedParamKey, uint256 time, byte[] extraData)
func (r *EvmRegistry) decodeComposerRequest(data []byte) (*ComposerRequestV1, error) {
	e := r.composer.abi.Errors["ComposerRequestV1"]
	fmt.Println("DATA", hexutil.Encode(data))
	fmt.Println(e)
	unpack, err := e.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	return &ComposerRequestV1{
		scriptHash:         *abi.ConvertType(errorParameters[0], new(string)).(*string),
		functionsArguments: *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		useMercury:         *abi.ConvertType(errorParameters[2], new(bool)).(*bool),
		feedParamKey:       *abi.ConvertType(errorParameters[3], new(string)).(*string),
		feeds:              *abi.ConvertType(errorParameters[4], new([]string)).(*[]string),
		timeParamKey:       *abi.ConvertType(errorParameters[5], new(string)).(*string),
		time:               *abi.ConvertType(errorParameters[6], new(*big.Int)).(**big.Int),
		extraData:          *abi.ConvertType(errorParameters[7], new([]byte)).(*[]byte),
	}, nil
}

func (r *EvmRegistry) checkCallback(ctx context.Context, values [][]byte, lookup *StreamsLookup) (encoding.PipelineExecutionState, bool, hexutil.Bytes, error) {
	payload, err := r.abi.Pack("checkCallback", lookup.upkeepId, values, lookup.extraData)
	if err != nil {
		return encoding.PackUnpackDecodeFailed, false, nil, err
	}

	var b hexutil.Bytes
	args := map[string]interface{}{
		"to":   r.addr.Hex(),
		"data": hexutil.Bytes(payload),
	}

	// call checkCallback function at the block which OCR3 has agreed upon
	err = r.client.CallContext(ctx, &b, "eth_call", args, hexutil.EncodeUint64(lookup.block))
	if err != nil {
		return encoding.RpcFlakyFailure, true, nil, err
	}
	return encoding.NoPipelineError, false, b, nil
}

// doMercuryRequest sends requests to Mercury API to retrieve mercury data.
func (r *EvmRegistry) doMercuryRequest(ctx context.Context, sl *StreamsLookup, lggr logger.Logger) (encoding.PipelineExecutionState, encoding.UpkeepFailureReason, [][]byte, bool, error) {
	var isMercuryV03 bool
	resultLen := len(sl.feeds)
	ch := make(chan MercuryData, resultLen)
	if len(sl.feeds) == 0 {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", sl.feedParamKey, sl.timeParamKey, sl.feeds)
	}
	if sl.feedParamKey == feedIdHex && sl.timeParamKey == blockNumber {
		// only mercury v0.2
		for i := range sl.feeds {
			go r.singleFeedRequest(ctx, ch, i, sl, lggr)
		}
	} else if sl.feedParamKey == feedIDs && sl.timeParamKey == timestamp {
		// only mercury v0.3
		resultLen = 1
		isMercuryV03 = true
		ch = make(chan MercuryData, resultLen)
		go r.multiFeedsRequest(ctx, ch, sl, lggr)
	} else {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonInvalidRevertDataInput, [][]byte{}, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", sl.feedParamKey, sl.timeParamKey, sl.feeds)
	}

	var reqErr error
	results := make([][]byte, len(sl.feeds))
	retryable := true
	allSuccess := true
	// in v0.2, use the last execution error as the state, if no execution errors, state will be no error
	state := encoding.NoPipelineError
	for i := 0; i < resultLen; i++ {
		m := <-ch
		if m.Error != nil {
			reqErr = errors.Join(reqErr, m.Error)
			retryable = retryable && m.Retryable
			allSuccess = false
			if m.State != encoding.NoPipelineError {
				state = m.State
			}
			continue
		}
		if isMercuryV03 {
			results = m.Bytes
		} else {
			results[m.Index] = m.Bytes[0]
		}
	}
	// only retry when not all successful AND none are not retryable
	return state, encoding.UpkeepFailureReasonNone, results, retryable && !allSuccess, reqErr
}

// singleFeedRequest sends a v0.2 Mercury request for a single feed report.
func (r *EvmRegistry) singleFeedRequest(ctx context.Context, ch chan<- MercuryData, index int, sl *StreamsLookup, lggr logger.Logger) {
	q := url.Values{
		sl.feedParamKey: {sl.feeds[index]},
		sl.timeParamKey: {sl.time.String()},
	}
	mercuryURL := r.mercury.cred.URL
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, mercuryPathV02, q.Encode())
	lggr.Debugf("request URL for upkeep %s feed %s: %s", sl.upkeepId.String(), sl.feeds[index], reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: index, Error: err, Retryable: false, State: encoding.InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, mercuryPathV02+q.Encode(), []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set(headerContentType, applicationJson)
	req.Header.Set(headerAuthorization, r.mercury.cred.Username)
	req.Header.Set(headerTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(headerSignature, signature)

	// in the case of multiple retries here, use the last attempt's data
	state := encoding.NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s GET request fails for feed %s: %v", sl.time.String(), sl.upkeepId.String(), sl.feeds[index], err1)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err1
			}
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					lggr.Warnf("failed to close mercury response Body: %s", err)
				}
			}(resp.Body)

			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = encoding.InvalidMercuryResponse
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				lggr.Warnf("at block %s upkeep %s received status code %d for feed %s", sl.time.String(), sl.upkeepId.String(), resp.StatusCode, sl.feeds[index])
				retryable = true
				state = encoding.MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at block %s upkeep %s received status code %d for feed %s", sl.time.String(), sl.upkeepId.String(), resp.StatusCode, sl.feeds[index])
			}

			lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.2 with BODY=%s", sl.time.String(), sl.upkeepId.String(), resp.StatusCode, hexutil.Encode(body))

			var m MercuryV02Response
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s failed to unmarshal body to MercuryV02Response for feed %s: %v", sl.time.String(), sl.upkeepId.String(), sl.feeds[index], err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				lggr.Warnf("at block %s upkeep %s failed to decode chainlinkBlob %s for feed %s: %v", sl.time.String(), sl.upkeepId.String(), m.ChainlinkBlob, sl.feeds[index], err1)
				retryable = false
				state = encoding.InvalidMercuryResponse
				return err1
			}
			ch <- MercuryData{
				Index:     index,
				Bytes:     [][]byte{blobBytes},
				Retryable: false,
				State:     encoding.NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(retryDelay),
		retry.Attempts(totalAttempt))

	if !sent {
		md := MercuryData{
			Index:     index,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     fmt.Errorf("failed to request feed for %s: %w", sl.feeds[index], retryErr),
			State:     state,
		}
		ch <- md
	}
}

// multiFeedsRequest sends a Mercury v0.3 request for a multi-feed report
func (r *EvmRegistry) multiFeedsRequest(ctx context.Context, ch chan<- MercuryData, sl *StreamsLookup, lggr logger.Logger) {
	// this won't work bc q.Encode() will encode commas as '%2C' but the server is strictly expecting a comma separated list
	//q := url.Values{
	//	feedIDs:   {strings.Join(sl.feeds, ",")},
	//	timestamp: {sl.time.String()},
	//}
	params := fmt.Sprintf("%s=%s&%s=%s", feedIDs, strings.Join(sl.feeds, ","), timestamp, sl.time.String())
	reqUrl := fmt.Sprintf("%s%s%s", r.mercury.cred.URL, mercuryBatchPathV03, params)
	lggr.Debugf("request URL for upkeep %s userId %s: %s", sl.upkeepId.String(), r.mercury.cred.Username, reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: 0, Error: err, Retryable: false, State: encoding.InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := r.generateHMAC(http.MethodGet, mercuryBatchPathV03+params, []byte{}, r.mercury.cred.Username, r.mercury.cred.Password, ts)
	req.Header.Set(headerContentType, applicationJson)
	// username here is often referred to as user id
	req.Header.Set(headerAuthorization, r.mercury.cred.Username)
	req.Header.Set(headerTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(headerSignature, signature)
	// mercury will inspect authorization headers above to make sure this user (in automation's context, this node) is eligible to access mercury
	// and if it has an automation role. it will then look at this upkeep id to check if it has access to all the requested feeds.
	req.Header.Set(headerUpkeepId, sl.upkeepId.String())

	// in the case of multiple retries here, use the last attempt's data
	state := encoding.NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := r.hc.Do(req)
			if err1 != nil {
				lggr.Warnf("at timestamp %s upkeep %s GET request fails from mercury v0.3: %v", sl.time.String(), sl.upkeepId.String(), err1)
				retryable = true
				state = encoding.MercuryFlakyFailure
				return err1
			}
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					lggr.Warnf("failed to close mercury response Body: %s", err)
				}
			}(resp.Body)

			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = encoding.InvalidMercuryResponse
				return err1
			}

			lggr.Infof("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.time.String(), sl.upkeepId.String(), resp.StatusCode)
			if resp.StatusCode == http.StatusUnauthorized {
				retryable = false
				state = encoding.UpkeepNotAuthorized
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by unauthorized upkeep", sl.time.String(), sl.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode == http.StatusBadRequest {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by invalid format of timestamp", sl.time.String(), sl.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode == http.StatusInternalServerError {
				retryable = true
				state = encoding.MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusInternalServerError)
			} else if resp.StatusCode == 420 {
				// in 0.3, this will happen when missing/malformed query args, missing or bad required headers, non-existent feeds, or no permissions for feeds
				retryable = false
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by missing/malformed query args, missing or bad required headers, non-existent feeds, or no permissions for feeds", sl.time.String(), sl.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = encoding.InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3", sl.time.String(), sl.upkeepId.String(), resp.StatusCode)
			}

			lggr.Debugf("at block %s upkeep %s received status code %d from mercury v0.3 with BODY=%s", sl.time.String(), sl.upkeepId.String(), resp.StatusCode, hexutil.Encode(body))

			var response MercuryV03Response
			err1 = json.Unmarshal(body, &response)
			if err1 != nil {
				lggr.Warnf("at timestamp %s upkeep %s failed to unmarshal body to MercuryV03Response from mercury v0.3: %v", sl.time.String(), sl.upkeepId.String(), err1)
				retryable = false
				state = encoding.MercuryUnmarshalError
				return err1
			}
			// in v0.3, if some feeds are not available, the server will only return available feeds, but we need to make sure ALL feeds are retrieved before calling user contract
			// hence, retry in this case. retry will help when we send a very new timestamp and reports are not yet generated
			if len(response.Reports) != len(sl.feeds) {
				// TODO: AUTO-5044: calculate what reports are missing and log a warning
				retryable = true
				state = encoding.MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusNotFound)
			}
			var reportBytes [][]byte
			for _, rsp := range response.Reports {
				reportBytes = append(reportBytes, rsp.FullReport)
			}
			ch <- MercuryData{
				Index:     0,
				Bytes:     reportBytes,
				Retryable: false,
				State:     encoding.NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(retryDelay),
		retry.Attempts(totalAttempt))

	if !sent {
		md := MercuryData{
			Index:     0,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     retryErr,
			State:     state,
		}
		ch <- md
	}
}

// generateHMAC calculates a user HMAC for Mercury server authentication.
func (r *EvmRegistry) generateHMAC(method string, path string, body []byte, clientId string, secret string, ts int64) string {
	bodyHash := sha256.New()
	bodyHash.Write(body)
	hashString := fmt.Sprintf("%s %s %s %s %d",
		method,
		path,
		hex.EncodeToString(bodyHash.Sum(nil)),
		clientId,
		ts)
	signedMessage := hmac.New(sha256.New, []byte(secret))
	signedMessage.Write([]byte(hashString))
	userHmac := hex.EncodeToString(signedMessage.Sum(nil))
	return userHmac
}
