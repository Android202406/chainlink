package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPrometheusExporter(t *testing.T) {
	t.Run("should set correct labels and cleanup all labels associated with different transmitters", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := new(MetricsMock)
		metrics.Test(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		metrics.On("SetFeedContractMetadata",
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
			feedConfig.GetSymbol(),         // symbol
		).Once()
		exporter, err := factory.NewExporter(chainConfig, feedConfig)
		require.NoError(t, err)

		envelope1, err := generateEnvelope()
		require.NoError(t, err)
		envelope1.Transmitter = types.Account(hexutil.Encode([]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, uint8(1),
		}))
		envelope2, err := generateEnvelope()
		require.NoError(t, err)
		envelope2.Transmitter = types.Account(hexutil.Encode([]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, uint8(2),
		}))

		humanizedAnswer1 := float64(envelope1.LatestAnswer.Uint64()) / float64(feedConfig.GetMultiply().Uint64())
		humanizedAnswer2 := float64(envelope2.LatestAnswer.Uint64()) / float64(feedConfig.GetMultiply().Uint64())
		humanizedJuelsPerFeeCoin1 := float64(envelope1.JuelsPerFeeCoin.Uint64()) / float64(feedConfig.GetMultiply().Uint64())
		humanizedJuelsPerFeeCoin2 := float64(envelope2.JuelsPerFeeCoin.Uint64()) / float64(feedConfig.GetMultiply().Uint64())

		metrics.On("SetFeedContractLinkBalance",
			float64(envelope1.LinkBalance.Uint64()), // balance
			feedConfig.GetID(),                      // contractAddress
			feedConfig.GetID(),                      // feedID
			chainConfig.GetChainID(),                // chainID
			feedConfig.GetContractStatus(),          // contractStatus
			feedConfig.GetContractType(),            // contractType
			feedConfig.GetName(),                    // feedName
			feedConfig.GetPath(),                    // feedPath
			chainConfig.GetNetworkID(),              // networkID
			chainConfig.GetNetworkName(),            // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(envelope1.Transmitter), // oracleName
			string(envelope1.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope1.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
		metrics.On("SetOffchainAggregatorAnswers",
			humanizedAnswer1,               // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorAnswersRaw",
			float64(envelope1.LatestAnswer.Uint64()), // answer
			feedConfig.GetID(),                       // contractAddress
			feedConfig.GetID(),                       // feedID
			chainConfig.GetChainID(),                 // chainID
			feedConfig.GetContractStatus(),           // contractStatus
			feedConfig.GetContractType(),             // contractType
			feedConfig.GetName(),                     // feedName
			feedConfig.GetPath(),                     // feedPath
			chainConfig.GetNetworkID(),               // networkID
			chainConfig.GetNetworkName(),             // networkName
		).Once()
		metrics.On("IncOffchainAggregatorAnswersTotal",
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinRaw",
			float64(envelope1.JuelsPerFeeCoin.Uint64()),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoin",
			humanizedJuelsPerFeeCoin1,
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorAnswerStalled",
			mock.Anything,                  // isSet
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			humanizedAnswer1,               // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			string(envelope1.Transmitter),  // sender
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorRoundID",
			float64(envelope1.AggregatorRoundID),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		exporter.Export(ctx, envelope1)

		metrics.On("SetFeedContractLinkBalance",
			float64(envelope2.LinkBalance.Uint64()), // balance
			feedConfig.GetID(),                      // contractAddress
			feedConfig.GetID(),                      // feedID
			chainConfig.GetChainID(),                // chainID
			feedConfig.GetContractStatus(),          // contractStatus
			feedConfig.GetContractType(),            // contractType
			feedConfig.GetName(),                    // feedName
			feedConfig.GetPath(),                    // feedPath
			chainConfig.GetNetworkID(),              // networkID
			chainConfig.GetNetworkName(),            // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(envelope2.Transmitter), // oracleName
			string(envelope2.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope2.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
		metrics.On("SetOffchainAggregatorAnswers",
			humanizedAnswer2,
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorAnswersRaw",
			float64(envelope2.LatestAnswer.Uint64()), // answer
			feedConfig.GetID(),                       // contractAddress
			feedConfig.GetID(),                       // feedID
			chainConfig.GetChainID(),                 // chainID
			feedConfig.GetContractStatus(),           // contractStatus
			feedConfig.GetContractType(),             // contractType
			feedConfig.GetName(),                     // feedName
			feedConfig.GetPath(),                     // feedPath
			chainConfig.GetNetworkID(),               // networkID
			chainConfig.GetNetworkName(),             // networkName
		).Once()
		metrics.On("IncOffchainAggregatorAnswersTotal",
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinRaw",
			float64(envelope2.JuelsPerFeeCoin.Uint64()),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoin",
			humanizedJuelsPerFeeCoin2,
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorAnswerStalled",
			mock.Anything,                  // isSet
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			humanizedAnswer2,
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			string(envelope2.Transmitter),  // sender
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorRoundID",
			float64(envelope2.AggregatorRoundID),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		exporter.Export(ctx, envelope2)

		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(envelope1.Transmitter),  // oracleName
			string(envelope1.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		).Once()
		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(envelope2.Transmitter),  // oracleName
			string(envelope2.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		).Once()
		exporter.Cleanup(ctx)

		mock.AssertExpectationsForObjects(t, metrics)
	})
	t.Run("should not emit metrics for stale transmissions", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := new(MetricsMock)
		metrics.Test(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		metrics.On("SetFeedContractMetadata",
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
			feedConfig.GetSymbol(),         // symbol
		).Once()
		exporter, err := factory.NewExporter(chainConfig, feedConfig)
		require.NoError(t, err)

		envelope1, err := generateEnvelope()
		require.NoError(t, err)
		envelope2, err := generateEnvelope()
		require.NoError(t, err)
		envelope2.LatestAnswer = envelope1.LatestAnswer
		envelope2.LatestTimestamp = envelope1.LatestTimestamp
		envelope2.Transmitter = envelope1.Transmitter

		humanizedAnswer := float64(envelope1.LatestAnswer.Uint64()) / float64(feedConfig.GetMultiply().Uint64())
		humanizedJuelsPerFeeCoin := float64(envelope1.JuelsPerFeeCoin.Uint64()) / float64(feedConfig.GetMultiply().Uint64())

		metrics.On("SetFeedContractLinkBalance",
			float64(envelope1.LinkBalance.Uint64()), // balance
			feedConfig.GetID(),                      // contractAddress
			feedConfig.GetID(),                      // feedID
			chainConfig.GetChainID(),                // chainID
			feedConfig.GetContractStatus(),          // contractStatus
			feedConfig.GetContractType(),            // contractType
			feedConfig.GetName(),                    // feedName
			feedConfig.GetPath(),                    // feedPath
			chainConfig.GetNetworkID(),              // networkID
			chainConfig.GetNetworkName(),            // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(envelope1.Transmitter), // oracleName
			string(envelope1.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope1.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
		metrics.On("SetOffchainAggregatorAnswers",
			humanizedAnswer,                // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorAnswersRaw",
			float64(envelope1.LatestAnswer.Uint64()), // answer
			feedConfig.GetID(),                       // contractAddress
			feedConfig.GetID(),                       // feedID
			chainConfig.GetChainID(),                 // chainID
			feedConfig.GetContractStatus(),           // contractStatus
			feedConfig.GetContractType(),             // contractType
			feedConfig.GetName(),                     // feedName
			feedConfig.GetPath(),                     // feedPath
			chainConfig.GetNetworkID(),               // networkID
			chainConfig.GetNetworkName(),             // networkName
		).Once()
		metrics.On("IncOffchainAggregatorAnswersTotal",
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinRaw",
			float64(envelope1.JuelsPerFeeCoin.Uint64()),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoin",
			humanizedJuelsPerFeeCoin,
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorAnswerStalled",
			mock.Anything,                  // isSet
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			humanizedAnswer,                // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			string(envelope1.Transmitter),  // sender
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorRoundID",
			float64(envelope1.AggregatorRoundID),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		exporter.Export(ctx, envelope1)

		metrics.On("SetFeedContractLinkBalance",
			float64(envelope2.LinkBalance.Uint64()), // balance
			feedConfig.GetID(),                      // contractAddress
			feedConfig.GetID(),                      // feedID
			chainConfig.GetChainID(),                // chainID
			feedConfig.GetContractStatus(),          // contractStatus
			feedConfig.GetContractType(),            // contractType
			feedConfig.GetName(),                    // feedName
			feedConfig.GetPath(),                    // feedPath
			chainConfig.GetNetworkID(),              // networkID
			chainConfig.GetNetworkName(),            // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(envelope2.Transmitter), // oracleName
			string(envelope2.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope2.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
		metrics.On("SetOffchainAggregatorAnswerStalled",
			mock.Anything,                  // isSet
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		exporter.Export(ctx, envelope2)

		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(envelope1.Transmitter),  // oracleName
			string(envelope1.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		).Once()
		exporter.Cleanup(ctx)

		metrics.AssertNumberOfCalls(t, "SetOffchainAggregatorAnswers", 1)
		metrics.AssertNumberOfCalls(t, "SetOffchainAggregatorAnswersRaw", 1)
		metrics.AssertNumberOfCalls(t, "IncOffchainAggregatorAnswersTotal", 1)
		metrics.AssertNumberOfCalls(t, "SetOffchainAggregatorSubmissionReceivedValues", 1)
		mock.AssertExpectationsForObjects(t, metrics)
	})
	t.Run("should emit transaction results metrics", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := new(MetricsMock)
		metrics.Test(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()

		metrics.On("SetFeedContractMetadata",
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
			feedConfig.GetSymbol(),         // symbol
		).Once()
		exporter, err := factory.NewExporter(chainConfig, feedConfig)
		require.NoError(t, err)

		txResults := generateTxResults()
		metrics.On("SetFeedContractTransmissionsSucceeded",
			float64(txResults.NumSucceeded), // succeeded
			feedConfig.GetID(),              // contractAddress
			feedConfig.GetID(),              // feedID
			chainConfig.GetChainID(),        // chainID
			feedConfig.GetContractStatus(),  // contractStatus
			feedConfig.GetContractType(),    // contractType
			feedConfig.GetName(),            // feedName
			feedConfig.GetPath(),            // feedPath
			chainConfig.GetNetworkID(),      // networkID
			chainConfig.GetNetworkName(),    // networkName
		).Once()
		metrics.On("SetFeedContractTransmissionsFailed",
			float64(txResults.NumFailed),   // failed
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		exporter.Export(ctx, txResults)

		mock.AssertExpectationsForObjects(t, metrics)
	})
}
