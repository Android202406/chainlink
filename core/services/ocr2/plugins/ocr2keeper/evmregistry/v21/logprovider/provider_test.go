package logprovider

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestLogEventProvider_GetFilters(t *testing.T) {
	p := NewLogProvider(logger.TestLogger(t), nil, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200, big.NewInt(1)))

	_, f := newEntry(p, 1)
	p.filterStore.AddActiveUpkeeps(f)

	t.Run("no filters", func(t *testing.T) {
		filters := p.getFilters(0, big.NewInt(0))
		require.Len(t, filters, 1)
		require.Equal(t, len(filters[0].addr), 0)
	})

	t.Run("has filter with lower lastPollBlock", func(t *testing.T) {
		filters := p.getFilters(0, f.upkeepID)
		require.Len(t, filters, 1)
		require.Greater(t, len(filters[0].addr), 0)
		filters = p.getFilters(10, f.upkeepID)
		require.Len(t, filters, 1)
		require.Greater(t, len(filters[0].addr), 0)
	})

	t.Run("has filter with higher lastPollBlock", func(t *testing.T) {
		_, f := newEntry(p, 2)
		f.lastPollBlock = 3
		p.filterStore.AddActiveUpkeeps(f)

		filters := p.getFilters(1, f.upkeepID)
		require.Len(t, filters, 1)
		require.Equal(t, len(filters[0].addr), 0)
	})

	t.Run("has filter with higher configUpdateBlock", func(t *testing.T) {
		_, f := newEntry(p, 2)
		f.configUpdateBlock = 3
		p.filterStore.AddActiveUpkeeps(f)

		filters := p.getFilters(1, f.upkeepID)
		require.Len(t, filters, 1)
		require.Equal(t, len(filters[0].addr), 0)
	})
}

func TestLogEventProvider_UpdateEntriesLastPoll(t *testing.T) {
	p := NewLogProvider(logger.TestLogger(t), nil, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), NewOptions(200, big.NewInt(1)))

	n := 10

	// entries := map[string]upkeepFilter{}
	for i := 0; i < n; i++ {
		_, f := newEntry(p, i+1)
		p.filterStore.AddActiveUpkeeps(f)
	}

	t.Run("no entries", func(t *testing.T) {
		_, f := newEntry(p, n*2)
		f.lastPollBlock = 10
		p.updateFiltersLastPoll([]upkeepFilter{f})

		filters := p.filterStore.GetFilters(nil)
		for _, f := range filters {
			require.Equal(t, int64(0), f.lastPollBlock)
		}
	})

	t.Run("update entries", func(t *testing.T) {
		_, f2 := newEntry(p, n-2)
		f2.lastPollBlock = 10
		_, f1 := newEntry(p, n-1)
		f1.lastPollBlock = 10
		p.updateFiltersLastPoll([]upkeepFilter{f1, f2})

		p.filterStore.RangeFiltersByIDs(func(_ int, f upkeepFilter) {
			require.Equal(t, int64(10), f.lastPollBlock)
		}, f1.upkeepID, f2.upkeepID)

		// update with same block
		p.updateFiltersLastPoll([]upkeepFilter{f1})

		// checking other entries are not updated
		_, f := newEntry(p, 1)
		p.filterStore.RangeFiltersByIDs(func(_ int, f upkeepFilter) {
			require.Equal(t, int64(0), f.lastPollBlock)
		}, f.upkeepID)
	})
}

func TestLogEventProvider_ScheduleReadJobs(t *testing.T) {
	mp := new(mocks.LogPoller)

	tests := []struct {
		name         string
		maxBatchSize int
		ids          []int
		addrs        []string
	}{
		{
			"no entries",
			3,
			[]int{},
			[]string{},
		},
		{
			"single entry",
			3,
			[]int{1},
			[]string{"0x1111111"},
		},
		{
			"happy flow",
			3,
			[]int{1, 2, 3},
			[]string{"0x1111111", "0x2222222", "0x3333333"},
		},
		{
			"batching",
			3,
			[]int{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10, 11, 12,
				13, 14, 15,
				16, 17, 18,
				19, 20, 21,
			},
			[]string{
				"0x11111111",
				"0x22222222",
				"0x33333333",
				"0x111111111",
				"0x122222222",
				"0x133333333",
				"0x1111111111",
				"0x1122222222",
				"0x1133333333",
				"0x11111111111",
				"0x11122222222",
				"0x11133333333",
				"0x111111111111",
				"0x111122222222",
				"0x111133333333",
				"0x1111111111111",
				"0x1111122222222",
				"0x1111133333333",
				"0x11111111111111",
				"0x11111122222222",
				"0x11111133333333",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)

			readInterval := 10 * time.Millisecond
			opts := NewOptions(200, big.NewInt(1))
			opts.ReadInterval = readInterval

			p := NewLogProvider(logger.TestLogger(t), mp, big.NewInt(1), &mockedPacker{}, NewUpkeepFilterStore(), opts)

			var ids []*big.Int
			for i, id := range tc.ids {
				_, f := newEntry(p, id, tc.addrs[i])
				p.filterStore.AddActiveUpkeeps(f)
				ids = append(ids, f.upkeepID)
			}

			reads := make(chan []*big.Int, 100)

			go func(ctx context.Context) {
				p.scheduleReadJobs(ctx, func(ids []*big.Int) {
					select {
					case reads <- ids:
					default:
						t.Log("dropped ids")
					}
				})
			}(ctx)

			batches := (len(tc.ids) / tc.maxBatchSize) + 1

			timeoutTicker := time.NewTicker(readInterval * time.Duration(batches*10))
			defer timeoutTicker.Stop()

			got := map[string]int{}

		readLoop:
			for {
				select {
				case <-timeoutTicker.C:
					break readLoop
				case batch := <-reads:
					for _, id := range batch {
						got[id.String()]++
					}
				case <-ctx.Done():
					break readLoop
				default:
					if p.CurrentPartitionIdx() > uint64(batches+1) {
						break readLoop
					}
				}
				runtime.Gosched()
			}

			require.Equal(t, len(ids), len(got))
			for _, id := range ids {
				_, ok := got[id.String()]
				require.True(t, ok, "id not found %s", id.String())
				require.GreaterOrEqual(t, got[id.String()], 1, "id don't have schdueled job %s", id.String())
			}
		})
	}
}

func TestLogEventProvider_ReadLogs(t *testing.T) {
	ctx := testutils.Context(t)

	mp := new(mocks.LogPoller)

	mp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
	mp.On("ReplayAsync", mock.Anything).Return()
	mp.On("HasFilter", mock.Anything).Return(false)
	mp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
	mp.On("LatestBlock", mock.Anything).Return(logpoller.LogPollerBlock{BlockNumber: int64(1)}, nil)
	mp.On("LogsWithSigs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{
		{
			BlockNumber: 1,
			TxHash:      common.HexToHash("0x1"),
		},
	}, nil)

	filterStore := NewUpkeepFilterStore()
	p := NewLogProvider(logger.TestLogger(t), mp, big.NewInt(1), &mockedPacker{}, filterStore, NewOptions(200, big.NewInt(1)))

	var ids []*big.Int
	for i := 0; i < 10; i++ {
		cfg, f := newEntry(p, i+1)
		ids = append(ids, f.upkeepID)
		require.NoError(t, p.RegisterFilter(ctx, FilterOptions{
			UpkeepID:      f.upkeepID,
			TriggerConfig: cfg,
		}))
	}

	t.Run("no entries", func(t *testing.T) {
		require.NoError(t, p.ReadLogs(ctx, big.NewInt(999999)))
		logs := p.buffer.peek(10)
		require.Len(t, logs, 0)
	})

	t.Run("has entries", func(t *testing.T) {
		require.NoError(t, p.ReadLogs(ctx, ids[:2]...))
		logs := p.buffer.peek(10)
		require.Len(t, logs, 2)

		var updatedFilters []upkeepFilter
		filterStore.RangeFiltersByIDs(func(i int, f upkeepFilter) {
			updatedFilters = append(updatedFilters, f.Clone())
		}, ids[:2]...)
		for _, f := range updatedFilters {
			// Last poll block should be updated
			require.Equal(t, int64(1), f.lastPollBlock)
		}
	})

	// TODO: test rate limiting

}

func newEntry(p *logEventProvider, i int, args ...string) (LogTriggerConfig, upkeepFilter) {
	idBytes := append(common.LeftPadBytes([]byte{1}, 16), []byte(fmt.Sprintf("%d", i))...)
	id := ocr2keepers.UpkeepIdentifier{}
	copy(id[:], idBytes)
	uid := id.BigInt()
	for len(args) < 2 {
		args = append(args, "0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d")
	}
	addr, topic0 := args[0], args[1]
	cfg := LogTriggerConfig{
		ContractAddress: common.HexToAddress(addr),
		FilterSelector:  0,
		Topic0:          common.HexToHash(topic0),
	}
	filter := p.newLogFilter(uid, cfg)
	topics := make([]common.Hash, len(filter.EventSigs))
	copy(topics, filter.EventSigs)
	f := upkeepFilter{
		upkeepID: uid,
		addr:     filter.Addresses[0].Bytes(),
		topics:   topics,
	}
	return cfg, f
}

func TestLogEventProvider_GetLatestPayloads(t *testing.T) {
	t.Run("small number of upkeeps, 100 logs per upkeep per block for 100 blocks", func(t *testing.T) {
		upkeepIDs := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
			big.NewInt(4),
			big.NewInt(5),
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 5, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		// each upkeep should have 100 logs * 100 blocks = 10000 logs
		assert.Equal(t, 10000, len(bufV1.queues["1"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["2"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["3"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["4"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["5"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 9980, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9980, len(bufV1.queues["2"].logs))
		assert.Equal(t, 9980, len(bufV1.queues["3"].logs))
		assert.Equal(t, 9980, len(bufV1.queues["4"].logs))
		assert.Equal(t, 9980, len(bufV1.queues["5"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across the 5 upkeeps
		assert.Equal(t, 9960, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9960, len(bufV1.queues["2"].logs))
		assert.Equal(t, 9960, len(bufV1.queues["3"].logs))
		assert.Equal(t, 9960, len(bufV1.queues["4"].logs))
		assert.Equal(t, 9960, len(bufV1.queues["5"].logs))
	})

	t.Run("large number of upkeeps", func(t *testing.T) {
		var upkeepIDs []*big.Int

		for i := int64(1); i <= 200; i++ {
			upkeepIDs = append(upkeepIDs, big.NewInt(i))
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 200, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		// each upkeep should have 100 logs * 100 blocks = 10000 logs
		assert.Equal(t, 10000, len(bufV1.queues["1"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["101"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["150"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.iterations)
		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps
		assert.Equal(t, 10000, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps
		assert.Equal(t, 9999, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps
		assert.Equal(t, 9999, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps
		assert.Equal(t, 9998, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["150"].logs))
	})

	t.Run("large number of upkeeps, with logs added mid way through reading", func(t *testing.T) {
		var upkeepIDs []*big.Int

		for i := int64(1); i <= 200; i++ {
			upkeepIDs = append(upkeepIDs, big.NewInt(i))
		}

		filterStore := NewUpkeepFilterStore()

		logGenerator := func(start, end int64) []logpoller.Log {
			var res []logpoller.Log
			for i := start; i < end; i++ {
				for j := 0; j < 100; j++ {
					res = append(res, logpoller.Log{
						LogIndex:    int64(j),
						BlockHash:   common.HexToHash(fmt.Sprintf("%d", i+1)),
						BlockNumber: i + 1,
					})
				}
			}
			return res
		}

		// use a log poller that will create logs for the queried block range
		logPoller := &mockLogPoller{
			LatestBlockFn: func(ctx context.Context) (int64, error) {
				return 100, nil
			},
			LogsWithSigsFn: func(ctx context.Context, start, end int64, eventSigs []common.Hash, address common.Address) ([]logpoller.Log, error) {
				return logGenerator(start, end), nil
			},
		}

		// prepare the filter store with upkeeps
		for _, upkeepID := range upkeepIDs {
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		opts := NewOptions(200, big.NewInt(1))
		opts.BufferVersion = "v1"

		provider := NewLogProvider(logger.TestLogger(t), logPoller, big.NewInt(1), &mockedPacker{}, filterStore, opts)

		ctx := context.Background()

		err := provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 200, provider.bufferV1.NumOfUpkeeps())

		bufV1 := provider.bufferV1.(*logBuffer)

		// each upkeep should have 100 logs * 100 blocks = 10000 logs
		assert.Equal(t, 10000, len(bufV1.queues["1"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["9"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["21"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["101"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["150"].logs))

		payloads, err := provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.iterations)
		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps; with 2 iterations this means even upkeep IDs are dequeued first
		assert.Equal(t, 10000, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["40"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["45"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["50"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps; on the second iteration, odd upkeep IDs are dequeued
		assert.Equal(t, 9999, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["99"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["100"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps; on the third iteration, even upkeep IDs are dequeued once again
		assert.Equal(t, 9999, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["150"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["160"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["170"].logs))

		for i := int64(201); i <= 300; i++ {
			upkeepIDs = append(upkeepIDs, big.NewInt(i))
		}

		for i := 200; i < len(upkeepIDs); i++ {
			upkeepID := upkeepIDs[i]
			filterStore.AddActiveUpkeeps(
				upkeepFilter{
					addr:     []byte(upkeepID.String()),
					upkeepID: upkeepID,
					topics: []common.Hash{
						common.HexToHash(upkeepID.String()),
					},
				},
			)
		}

		err = provider.ReadLogs(ctx, upkeepIDs...)
		assert.NoError(t, err)

		assert.Equal(t, 300, provider.bufferV1.NumOfUpkeeps())

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 2, provider.iterations)
		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps; the new iterations
		// have not yet been recalculated despite the new logs being added; new iterations
		// are only calculated when current iteration maxes out at the total number of iterations
		assert.Equal(t, 9998, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["51"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["52"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		// with the newly added logs, iterations is recalculated
		assert.Equal(t, 3, provider.iterations)
		assert.Equal(t, 1, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps
		assert.Equal(t, 9998, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["11"].logs))
		assert.Equal(t, 9997, len(bufV1.queues["111"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9997, len(bufV1.queues["150"].logs))

		payloads, err = provider.GetLatestPayloads(ctx)
		assert.NoError(t, err)

		assert.Equal(t, 3, provider.iterations)
		assert.Equal(t, 2, provider.currentIteration)

		// we dequeue a maximum of 100 logs
		assert.Equal(t, 100, len(payloads))

		// the dequeue is evenly distributed across all upkeeps
		assert.Equal(t, 9997, len(bufV1.queues["1"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["2"].logs))
		assert.Equal(t, 9997, len(bufV1.queues["3"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["50"].logs))
		assert.Equal(t, 9998, len(bufV1.queues["101"].logs))
		assert.Equal(t, 9997, len(bufV1.queues["150"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["250"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["296"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["297"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["298"].logs))
		assert.Equal(t, 10000, len(bufV1.queues["299"].logs))
		assert.Equal(t, 9999, len(bufV1.queues["300"].logs))
	})
}

type mockedPacker struct {
}

func (p *mockedPacker) PackLogData(log logpoller.Log) ([]byte, error) {
	return log.Data, nil
}
