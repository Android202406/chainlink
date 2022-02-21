package monitoring

import (
	"context"
	"sync"
	"time"
)

type FeedMonitor interface {
	Run(ctx context.Context)
}

func NewFeedMonitor(
	log Logger,
	pollers []Poller,
	exporters []Exporter,
) FeedMonitor {
	return &feedMonitor{
		log,
		pollers,
		exporters,
	}
}

type feedMonitor struct {
	log       Logger
	pollers   []Poller
	exporters []Exporter
}

// Run should be executed as a goroutine.
// Signal termination by cancelling ctx; then wait for Run() to exit.
func (f *feedMonitor) Run(ctx context.Context) {
	f.log.Infow("starting feed monitor")
	wg := &sync.WaitGroup{}

	// Listen for updates
	updatesFanIn := make(chan interface{})
	wg.Add(len(f.pollers))
	for _, poller := range f.pollers {
		go func(poller Poller) {
			defer wg.Done()
			for {
				select {
				case update := <-poller.Updates():
					select {
					case updatesFanIn <- update:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}(poller)
	}

	// Consume updates.
CONSUME_LOOP:
	for {
		var update interface{}
		select {
		case update = <-updatesFanIn:
		case <-ctx.Done():
			break CONSUME_LOOP
		}
		// TODO (dru) do we need a worker pool here?
		wg.Add(len(f.exporters))
		for index, exp := range f.exporters {
			go func(index int, exp Exporter) {
				defer wg.Done()
				defer func() {
					if err := recover(); err != nil {
						f.log.Errorw("failed Export", "error", err, "index", index)
					}
				}()
				exp.Export(ctx, update)
			}(index, exp)
		}
	}

	// Cleanup happens after all the exporters have finished.
	wg.Wait()
	wg = &sync.WaitGroup{}
	defer wg.Wait()
	cleanupContext, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	wg.Add(len(f.exporters))
	for index, exp := range f.exporters {
		go func(index int, exp Exporter) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					f.log.Errorw("failed Cleanup", "error", err, "index", index)
				}
			}()
			exp.Cleanup(cleanupContext)
		}(index, exp)
	}
}
