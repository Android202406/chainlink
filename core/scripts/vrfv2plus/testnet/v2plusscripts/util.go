package v2plusscripts

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_sub_owner"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_consumer_example"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type VRFLoadTestMetrics struct {
	RequestCount                 *big.Int
	FulfilmentCount              *big.Int
	AverageFulfillmentInMillions *big.Int
	SlowestFulfillment           *big.Int
	FastestFulfillment           *big.Int
}

func DeployBHS(e helpers.Environment) (blockhashStoreAddress common.Address) {
	_, tx, _, err := blockhash_store.DeployBlockhashStore(e.Owner, e.Ec)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func DeployBatchBHS(e helpers.Environment, bhsAddress common.Address) (batchBHSAddress common.Address) {
	_, tx, _, err := batch_blockhash_store.DeployBatchBlockhashStore(e.Owner, e.Ec, bhsAddress)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func DeployCoordinator(
	e helpers.Environment,
	linkAddress string,
	bhsAddress string,
	linkEthAddress string,
) (coordinatorAddress common.Address) {
	_, tx, _, err := vrf_coordinator_v2_5.DeployVRFCoordinatorV25(
		e.Owner,
		e.Ec,
		common.HexToAddress(bhsAddress))
	helpers.PanicErr(err)
	coordinatorAddress = helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)

	// Set LINK and LINK ETH
	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(coordinatorAddress, e.Ec)
	helpers.PanicErr(err)

	linkTx, err := coordinator.SetLINKAndLINKNativeFeed(e.Owner,
		common.HexToAddress(linkAddress), common.HexToAddress(linkEthAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, linkTx, e.ChainID)
	return coordinatorAddress
}

func DeployBatchCoordinatorV2(e helpers.Environment, coordinatorAddress common.Address) (batchCoordinatorAddress common.Address) {
	_, tx, _, err := batch_vrf_coordinator_v2plus.DeployBatchVRFCoordinatorV2Plus(e.Owner, e.Ec, coordinatorAddress)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func EoaAddConsumerToSub(
	e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25,
	subID *big.Int,
	consumerAddress string,
) {
	txadd, err := coordinator.AddConsumer(e.Owner, subID, common.HexToAddress(consumerAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, txadd, e.ChainID)
}

func EoaCreateSub(e helpers.Environment, coordinator vrf_coordinator_v2_5.VRFCoordinatorV25) {
	tx, err := coordinator.CreateSubscription(e.Owner)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

// returns subscription ID that belongs to the given owner. Returns result found first
func FindSubscriptionID(e helpers.Environment, coordinator *vrf_coordinator_v2_5.VRFCoordinatorV25) *big.Int {
	// Use most recent 500 blocks as search window.
	head, err := e.Ec.BlockNumber(context.Background())
	helpers.PanicErr(err)
	fopts := &bind.FilterOpts{
		Start: head - 500,
	}

	subscriptionIterator, err := coordinator.FilterSubscriptionCreated(fopts, nil)
	helpers.PanicErr(err)

	if !subscriptionIterator.Next() {
		helpers.PanicErr(fmt.Errorf("expected at least 1 subID for the given owner %s", e.Owner.From.Hex()))
	}
	return subscriptionIterator.Event.SubId
}

func EoaDeployConsumer(e helpers.Environment,
	coordinatorAddress string,
	linkAddress string) (
	consumerAddress common.Address) {
	_, tx, _, err := vrf_v2plus_sub_owner.DeployVRFV2PlusExternalSubOwnerExample(
		e.Owner,
		e.Ec,
		common.HexToAddress(coordinatorAddress),
		common.HexToAddress(linkAddress))
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func EoaFundSubWithLink(
	e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25,
	linkAddress string, amount,
	subID *big.Int,
) {
	linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(linkAddress), e.Ec)
	helpers.PanicErr(err)
	bal, err := linkToken.BalanceOf(nil, e.Owner.From)
	helpers.PanicErr(err)
	fmt.Println("Initial account balance (Juels):", bal, e.Owner.From.String(), "Funding amount:", amount.String())
	b, err := utils.ABIEncode(`[{"type":"uint256"}]`, subID)
	helpers.PanicErr(err)
	tx, err := linkToken.TransferAndCall(e.Owner, coordinator.Address(), amount, b)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("sub ID: %d", subID))
}

func EoaFundSubWithNative(e helpers.Environment, coordinatorAddress common.Address, subID *big.Int, amount *big.Int) {
	coordinator, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(coordinatorAddress, e.Ec)
	helpers.PanicErr(err)
	e.Owner.Value = amount
	tx, err := coordinator.FundSubscriptionWithNative(e.Owner, subID)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func PrintCoordinatorConfig(coordinator *vrf_coordinator_v2_5.VRFCoordinatorV25) {
	cfg, err := coordinator.SConfig(nil)
	helpers.PanicErr(err)

	feeConfig, err := coordinator.SFeeConfig(nil)
	helpers.PanicErr(err)

	fmt.Printf("Coordinator config: %+v\n", cfg)
	fmt.Printf("Coordinator fee config: %+v\n", feeConfig)
}

func SetCoordinatorConfig(
	e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25,
	minConfs uint16,
	maxGasLimit uint32,
	stalenessSeconds uint32,
	gasAfterPayment uint32,
	fallbackWeiPerUnitLink *big.Int,
	feeConfig vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig,
) {
	tx, err := coordinator.SetConfig(
		e.Owner,
		minConfs,               // minRequestConfirmations
		maxGasLimit,            // max gas limit
		stalenessSeconds,       // stalenessSeconds
		gasAfterPayment,        // gasAfterPaymentCalculation
		fallbackWeiPerUnitLink, // 0.01 eth per link fallbackLinkPrice
		feeConfig,
	)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func RegisterCoordinatorProvingKey(e helpers.Environment,
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25, uncompressed string, oracleAddress string) {
	pubBytes, err := hex.DecodeString(uncompressed)
	helpers.PanicErr(err)
	pk, err := crypto.UnmarshalPubkey(pubBytes)
	helpers.PanicErr(err)
	tx, err := coordinator.RegisterProvingKey(e.Owner,
		common.HexToAddress(oracleAddress),
		[2]*big.Int{pk.X, pk.Y})
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(
		context.Background(),
		e.Ec,
		tx,
		e.ChainID,
		fmt.Sprintf("Uncompressed public key: %s,", uncompressed),
		fmt.Sprintf("Oracle address: %s,", oracleAddress),
	)
}

func WrapperDeploy(
	e helpers.Environment,
	link, linkEthFeed, coordinator common.Address,
) (common.Address, *big.Int) {
	address, tx, _, err := vrfv2plus_wrapper.DeployVRFV2PlusWrapper(e.Owner, e.Ec,
		link,
		linkEthFeed,
		coordinator)
	helpers.PanicErr(err)

	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	fmt.Println("VRFV2Wrapper address:", address)

	wrapper, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(address, e.Ec)
	helpers.PanicErr(err)

	subID, err := wrapper.SUBSCRIPTIONID(nil)
	helpers.PanicErr(err)
	fmt.Println("VRFV2Wrapper subscription id:", subID)

	return address, subID
}

func WrapperConfigure(
	e helpers.Environment,
	wrapperAddress common.Address,
	wrapperGasOverhead, coordinatorGasOverhead, premiumPercentage uint,
	keyHash string,
	maxNumWords uint,
	fallbackWeiPerUnitLink *big.Int,
	stalenessSeconds uint32,
	fulfillmentFlatFeeLinkPPM uint32,
	fulfillmentFlatFeeNativePPM uint32,
) {
	wrapper, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(wrapperAddress, e.Ec)
	helpers.PanicErr(err)

	tx, err := wrapper.SetConfig(
		e.Owner,
		uint32(wrapperGasOverhead),
		uint32(coordinatorGasOverhead),
		uint8(premiumPercentage),
		common.HexToHash(keyHash),
		uint8(maxNumWords),
		stalenessSeconds,
		fallbackWeiPerUnitLink,
		fulfillmentFlatFeeLinkPPM,
		fulfillmentFlatFeeNativePPM,
	)

	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func WrapperConsumerDeploy(
	e helpers.Environment,
	link, wrapper common.Address,
) common.Address {
	address, tx, _, err := vrfv2plus_wrapper_consumer_example.DeployVRFV2PlusWrapperConsumerExample(e.Owner, e.Ec,
		link,
		wrapper)
	helpers.PanicErr(err)

	helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
	fmt.Printf("VRFV2WrapperConsumerExample address: %s\n", address)
	return address
}

func EoaV2PlusLoadTestConsumerWithMetricsDeploy(e helpers.Environment, consumerCoordinator string) (consumerAddress common.Address) {
	_, tx, _, err := vrf_v2plus_load_test_with_metrics.DeployVRFV2PlusLoadTestWithMetrics(
		e.Owner,
		e.Ec,
		common.HexToAddress(consumerCoordinator),
	)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func LoadTestRequestRandomness(
	e helpers.Environment,
	consumerAddress common.Address,
	subID *big.Int,
	requestConfirmations uint,
	keyHashBytes common.Hash,
	cbGasLimit uint,
	nativePaymentEnabled bool,
	numWords uint,
	requests uint,
	runs uint,
) {
	consumer, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(
		consumerAddress,
		e.Ec)
	helpers.PanicErr(err)
	var txes []*types.Transaction
	for i := 0; i < int(runs); i++ {

		//todo - debug
		fmt.Println("consumer", consumer.Address().String())
		fmt.Println("e.Owner", e.Owner)
		fmt.Println("subID", subID)
		fmt.Println("requestConfirmations", requestConfirmations)
		fmt.Println("keyHashBytes", keyHashBytes)
		fmt.Println("cbGasLimit", cbGasLimit)
		fmt.Println("nativePaymentEnabled", nativePaymentEnabled)
		fmt.Println("numWords", numWords)
		fmt.Println("requests", requests)

		//todo - fails with "execution reverted" when calling this function from setup-env script,
		//even though when I run the same command in the console it works (go run . eoa-load-test-request-with-metrics )
		tx, err := consumer.RequestRandomWords(
			e.Owner,
			subID,
			uint16(requestConfirmations),
			keyHashBytes,
			uint32(cbGasLimit),
			nativePaymentEnabled,
			uint32(numWords),
			uint16(requests),
		)
		fmt.Println("tx", tx)
		helpers.PanicErr(err)
		fmt.Printf("TX %d: %s\n", i+1, helpers.ExplorerLink(e.ChainID, tx.Hash()))
		txes = append(txes, tx)
	}
	fmt.Println("Total number of requests sent:", (requests)*(runs))
	fmt.Println("fetching receipts for all transactions")
	for i, tx := range txes {
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("load test %d", i+1))
	}
}

func WaitForRequestCountEqualToFulfilmentCount(e helpers.Environment, consumerAddress common.Address, timeout time.Duration, wg *sync.WaitGroup) (*big.Int, *big.Int, error) {
	metricsChannel := make(chan *VRFLoadTestMetrics)
	metricsErrorChannel := make(chan error)

	testContext, testCancel := context.WithTimeout(context.Background(), timeout)
	defer testCancel()

	ticker := time.NewTicker(time.Second * 1)
	var metrics *VRFLoadTestMetrics
	for {
		select {
		case <-testContext.Done():
			ticker.Stop()
			wg.Done()
			return metrics.RequestCount, metrics.FulfilmentCount,
				fmt.Errorf("timeout waiting for rand request and fulfilments to be equal AFTER performance test was executed. Request Count: %d, Fulfilment Count: %d",
					metrics.RequestCount.Uint64(), metrics.FulfilmentCount.Uint64())
		case <-ticker.C:
			go getLoadTestMetrics(e, consumerAddress, metricsChannel)
		case metrics = <-metricsChannel:
			if metrics.RequestCount.Cmp(metrics.FulfilmentCount) == 0 {
				wg.Done()
				return metrics.RequestCount, metrics.FulfilmentCount, nil
			}
		case err := <-metricsErrorChannel:
			wg.Done()
			return nil, nil, err
		}
	}
}

func getLoadTestMetrics(
	e helpers.Environment,
	consumerAddress common.Address,
	metricsChannel chan *VRFLoadTestMetrics,
) {
	metrics := GetLoadTestMetricsFromConsumer(consumerAddress, e)
	metricsChannel <- metrics
}

func GetLoadTestMetricsFromConsumer(consumerAddress common.Address, e helpers.Environment) *VRFLoadTestMetrics {
	consumer, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(
		consumerAddress,
		e.Ec)
	helpers.PanicErr(err)
	responseCount, err := consumer.SResponseCount(nil)
	helpers.PanicErr(err)
	fmt.Println("Response Count: ", responseCount)
	requestCount, err := consumer.SRequestCount(nil)
	helpers.PanicErr(err)
	fmt.Println("Request Count: ", requestCount)
	averageFulfillmentInMillions, err := consumer.SAverageFulfillmentInMillions(nil)
	helpers.PanicErr(err)
	fmt.Println("Average Fulfillment In Millions: ", averageFulfillmentInMillions)
	slowestFulfillment, err := consumer.SSlowestFulfillment(nil)
	helpers.PanicErr(err)
	fmt.Println("Slowest Fulfillment: ", slowestFulfillment)
	fastestFulfillment, err := consumer.SFastestFulfillment(nil)
	helpers.PanicErr(err)
	fmt.Println("Fastest Fulfillment: ", fastestFulfillment)

	return &VRFLoadTestMetrics{
		requestCount,
		responseCount,
		averageFulfillmentInMillions,
		slowestFulfillment,
		fastestFulfillment,
	}
}
