package evm_test

//go:generate ./testfiles/chainlink_reader_test_setup.sh

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/testfiles"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const commonGasLimitOnEvms = uint64(4712388)

func TestChainReader(t *testing.T) {
	RunChainReaderInterfaceTests(t, &chainReaderInterfaceTester{})
}

type chainReaderInterfaceTester struct {
	chain       *mocks.Chain
	address     string
	address2    string
	chainConfig types.ChainReaderConfig
	auth        *bind.TransactOpts
	sim         *backends.SimulatedBackend
	pk          *ecdsa.PrivateKey
	evmTest     *testfiles.Testfiles
	cr          evm.ChainReaderService
}

func (it *chainReaderInterfaceTester) Setup(ctx context.Context, t *testing.T) {
	t.Cleanup(func() {
		// DB may be closed by the test already, ignore errors
		_ = it.cr.Close()
		it.cr = nil
		it.evmTest = nil
	})

	// can re-use the same chain for tests, just make new contract for each test
	if it.chain != nil {
		it.deployNewContracts(t)
		return
	}

	it.chain = &mocks.Chain{}
	it.setupChainNoClient(t)

	testStruct := CreateTestStruct(0, it)

	it.chainConfig = types.ChainReaderConfig{
		ChainContractReaders: map[string]types.ChainContractReader{
			AnyContractName: {
				ContractABI: testfiles.TestfilesMetaData.ABI,
				ChainReaderDefinitions: map[string]types.ChainReaderDefinition{
					MethodTakingLatestParamsReturningTestStruct: {
						ChainSpecificName: "GetElementAtIndex",
					},
					MethodReturningUint64: {
						ChainSpecificName: "GetPrimitiveValue",
					},
					DifferentMethodReturningUint64: {
						ChainSpecificName: "GetDifferentPrimitiveValue",
					},
					MethodReturningUint64Slice: {
						ChainSpecificName: "GetSliceValue",
					},
					EventName: {
						ChainSpecificName: "Triggered",
						ReadType:          types.Event,
					},
					MethodReturningSeenStruct: {
						ChainSpecificName: "ReturnSeen",
						InputModifications: codec.ModifiersConfig{
							&codec.HardCodeConfig{
								OnChainValues: map[string]any{
									"BigField": testStruct.BigField.String(),
									"Account":  hexutil.Encode(testStruct.Account),
								},
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.HardCodeConfig{
								OffChainValues: map[string]any{"ExtraField": anyExtraValue}},
						},
					},
				},
			},
			AnySecondContractName: {
				ContractABI: testfiles.TestfilesMetaData.ABI,
				ChainReaderDefinitions: map[string]types.ChainReaderDefinition{
					MethodReturningUint64: {
						ChainSpecificName: "GetDifferentPrimitiveValue",
					},
				},
			},
		},
	}
	it.chain.On("Client").Return(client.NewSimulatedBackendClient(t, it.sim, big.NewInt(1337)))
	it.deployNewContracts(t)
}

func (it *chainReaderInterfaceTester) Name() string {
	return "EVM"
}

func (it *chainReaderInterfaceTester) GetAccountBytes(i int) []byte {
	account := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	account[i%20] += byte(i)
	account[(i+3)%20] += byte(i + 3)
	return account[:]
}

func (it *chainReaderInterfaceTester) GetChainReader(ctx context.Context, t *testing.T) clcommontypes.ChainReader {
	if it.cr != nil {
		return it.cr
	}

	addr := common.HexToAddress(it.address)
	addr2 := common.HexToAddress(it.address2)
	lggr := logger.NullLogger
	db := pgtest.NewSqlxDB(t)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, db, lggr, pgtest.NewQConfig(true)), it.chain.Client(), lggr, time.Millisecond, false, 0, 1, 1, 10000)
	require.NoError(t, lp.Start(ctx))
	it.chain.On("LogPoller").Return(lp)
	b := evm.Bindings{
		AnyContractName: {
			MethodTakingLatestParamsReturningTestStruct: evm.NewAddrEvtFromAddress(addr),
			MethodReturningUint64:                       evm.NewAddrEvtFromAddress(addr),
			DifferentMethodReturningUint64:              evm.NewAddrEvtFromAddress(addr2),
			MethodReturningUint64Slice:                  evm.NewAddrEvtFromAddress(addr),
			EventName:                                   evm.NewAddrEvtFromAddress(addr),
			MethodReturningSeenStruct:                   evm.NewAddrEvtFromAddress(addr),
		},
		AnySecondContractName: {
			MethodReturningUint64: evm.NewAddrEvtFromAddress(addr2),
		},
	}
	cr, err := evm.NewChainReaderService(lggr, lp, b, it.chain, it.chainConfig)
	require.NoError(t, err)
	require.NoError(t, cr.Start(ctx))
	it.cr = cr
	return cr
}

func (it *chainReaderInterfaceTester) SetLatestValue(t *testing.T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*testfiles.TestfilesTransactor).AddTestStruct)
}

func (it *chainReaderInterfaceTester) TriggerEvent(t *testing.T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*testfiles.TestfilesTransactor).TriggerEvent)
}

type testStructFn = func(*testfiles.TestfilesTransactor, *bind.TransactOpts, int32, string, uint8, [32]uint8, common.Address, []common.Address, *big.Int, testfiles.MidLevelTestStruct) (*evmtypes.Transaction, error)

func (it *chainReaderInterfaceTester) sendTxWithTestStruct(t *testing.T, testStruct *TestStruct, fn testStructFn) {
	tx, err := fn(
		&it.evmTest.TestfilesTransactor,
		it.auth,
		testStruct.Field,
		testStruct.DifferentField,
		uint8(testStruct.OracleID),
		convertOracleIDs(testStruct.OracleIDs),
		common.Address(testStruct.Account),
		convertAccounts(testStruct.Accounts),
		testStruct.BigField,
		midToInternalType(testStruct.NestedStruct),
	)
	require.NoError(t, err)
	it.sim.Commit()
	it.incNonce()
	it.awaitTx(ctx, t, tx)
}

func convertOracleIDs(oracleIDs [32]commontypes.OracleID) [32]byte {
	convertedIds := [32]byte{}
	for i, id := range oracleIDs {
		convertedIds[i] = byte(id)
	}
	return convertedIds
}

func convertAccounts(accounts [][]byte) []common.Address {
	convertedAccounts := make([]common.Address, len(accounts))
	for i, a := range accounts {
		convertedAccounts[i] = common.Address(a)
	}
	return convertedAccounts
}

func (it *chainReaderInterfaceTester) setupChainNoClient(t require.TestingT) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	it.pk = privateKey

	it.auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	it.sim = backends.NewSimulatedBackend(core.GenesisAlloc{it.auth.From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
	it.sim.Commit()
}

func (it *chainReaderInterfaceTester) deployNewContracts(t *testing.T) {
	it.address = it.deployNewContract(t)
	it.address2 = it.deployNewContract(t)
}

func (it *chainReaderInterfaceTester) deployNewContract(t *testing.T) string {
	ctx := testutils.Context(t)
	gasPrice, err := it.sim.SuggestGasPrice(ctx)
	require.NoError(t, err)
	it.auth.GasPrice = gasPrice

	// 105528 was in the error: gas too low: have 0, want 105528
	// Not sure if there's a better way to get it.
	it.auth.GasLimit = 10552800

	address, tx, ts, err := testfiles.DeployTestfiles(it.auth, it.sim)

	require.NoError(t, err)
	it.sim.Commit()
	if it.evmTest == nil {
		it.evmTest = ts
	}
	it.incNonce()
	it.awaitTx(t, tx)
	return address.String()
}

func (it *chainReaderInterfaceTester) awaitTx(ctx context.Context, t *testing.T, tx *evmtypes.Transaction) {
	receipt, err := it.sim.TransactionReceipt(ctx, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, evmtypes.ReceiptStatusSuccessful, receipt.Status)
}

func (it *chainReaderInterfaceTester) incNonce() {
	if it.auth.Nonce == nil {
		it.auth.Nonce = big.NewInt(1)
	} else {
		it.auth.Nonce = it.auth.Nonce.Add(it.auth.Nonce, big.NewInt(1))
	}
}

func getAccounts(first TestStruct) []common.Address {
	accountBytes := make([]common.Address, len(first.Accounts))
	for i, account := range first.Accounts {
		accountBytes[i] = common.Address(account)
	}
	return accountBytes
}

func argsFromTestStruct(ts TestStruct) []any {
	return []any{
		ts.Field,
		ts.DifferentField,
		uint8(ts.OracleID),
		getOracleIDs(ts),
		common.Address(ts.Account),
		getAccounts(ts),
		ts.BigField,
		midToInternalType(ts.NestedStruct),
	}
}

func getOracleIDs(first TestStruct) [32]byte {
	oracleIDs := [32]byte{}
	for i, oracleID := range first.OracleIDs {
		oracleIDs[i] = byte(oracleID)
	}
	return oracleIDs
}

func toInternalType(testStruct TestStruct) testfiles.TestStruct {
	return testfiles.TestStruct{
		Field:          testStruct.Field,
		DifferentField: testStruct.DifferentField,
		OracleId:       byte(testStruct.OracleID),
		OracleIds:      convertOracleIDs(testStruct.OracleIDs),
		Account:        common.Address(testStruct.Account),
		Accounts:       convertAccounts(testStruct.Accounts),
		BigField:       testStruct.BigField,
		NestedStruct:   midToInternalType(testStruct.NestedStruct),
	}
}

func midToInternalType(m MidLevelTestStruct) testfiles.MidLevelTestStruct {
	return testfiles.MidLevelTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: testfiles.InnerTestStruct{
			I: int64(m.Inner.I),
			S: m.Inner.S,
		},
	}
}
