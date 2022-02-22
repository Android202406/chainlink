package monitoring

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const numPollerUpdates = 10
const numGoroutinesPerManaged = 10

func TestManager(t *testing.T) {
	t.Run("all goroutines are stopped before the new ones begin", func(t *testing.T) {
		// Poller fires 10 rounds of updates.
		// The manager identifies these updates, terminates the current running managed function and starts a new one.
		// The managed function in turn runs 10 noop goroutines and increments/decrements a goroutine counter.
		defer goleak.VerifyNone(t)

		var goRoutineCounter int64 = 0
		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		poller := &fakePoller{
			numPollerUpdates,
			make(chan interface{}),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			poller.Run(ctx)
		}()

		manager := NewManager(
			newNullLogger(),
			poller,
		)
		managed := func(ctx context.Context, _ []FeedConfig) {
			localWg := &sync.WaitGroup{}
			defer localWg.Wait()
			localWg.Add(numGoroutinesPerManaged)
			for i := 0; i < numGoroutinesPerManaged; i++ {
				go func(i int, ctx context.Context) {
					defer localWg.Done()
					atomic.AddInt64(&goRoutineCounter, 1)
					<-ctx.Done()
					atomic.AddInt64(&goRoutineCounter, -1)
				}(i, ctx)
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			manager.Run(ctx, managed)
		}()

		wg.Wait()
		require.Equal(t, int64(0), goRoutineCounter, "all child goroutines are gone")
	})

	t.Run("should not restart the monitor if the feeds are the same", func(t *testing.T) {
		feeds := []FeedConfig{
			generateFeedConfig(),
			generateFeedConfig(),
		}
		rddPoller := &fakePoller{0, make(chan interface{})}
		manager := NewManager(
			newNullLogger(),
			rddPoller,
		)

		var countManagedFuncExecutions uint64 = 0
		var managedFunc = func(_ context.Context, _ []FeedConfig) {
			atomic.AddUint64(&countManagedFuncExecutions, 1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			manager.Run(ctx, managedFunc)
		}()

		// The rdd poller returns the same feed configs three times!
		for i := 0; i < 3; i++ {
			select {
			case rddPoller.ch <- feeds:
			case <-ctx.Done():
			}
		}

		cancel()
		wg.Wait()

		require.Equal(t, countManagedFuncExecutions, uint64(1))
	})

	t.Run("should expose the current feeds to http", func(t *testing.T) {
		feeds := []FeedConfig{generateFeedConfig()}
		manager := &managerImpl{
			newNullLogger(),
			&fakePoller{0, make(chan interface{})},
			feeds,
			sync.Mutex{},
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/debug", nil)
		manager.HTTPHandler().ServeHTTP(rec, req)
		dec := json.NewDecoder(rec.Body)
		decodedFeeds := []fakeFeedConfig{}
		err := dec.Decode(&decodedFeeds)
		require.NoError(t, err)
		require.Equal(t, len(feeds), len(decodedFeeds))
		for i, feed := range feeds {
			require.Equal(t, feed, decodedFeeds[i])
		}
	})
}
