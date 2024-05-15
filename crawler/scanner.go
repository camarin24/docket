package docket

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type Scanner struct {
	queue         []time.Time
	lastExecution *time.Time
	interval      time.Duration
	mu            sync.Mutex
	logger        *zap.Logger
}

func (sc *Scanner) StartScheduler() {
	//TODO: Figure out how to properly stop this
	sc.logger.Info("Starting scheduler")
	ticker := time.NewTicker(sc.interval)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				sc.mu.Lock()
				sc.logger.Info("New task queued at", zap.Time("time", t))
				sc.queue = append(sc.queue, t)
				sc.mu.Unlock()
			}
		}
	}()
}

func (sc *Scanner) Next() *time.Time {
	return sc.Dequeue()
}

func (sc *Scanner) Dequeue() *time.Time {
	sc.logger.Info("Dequeueing task")
	var schedule time.Time
	for {
		// We do just to reduce the for oterations, i think
		//TODO: Search for a better alternative
		//time.Sleep(sc.interval / 5)
		sc.mu.Lock()
		if len(sc.queue) > 0 {
			schedule, sc.queue = sc.queue[len(sc.queue)-1], sc.queue[:len(sc.queue)-1]
			sc.mu.Unlock()
			break
		}
		sc.mu.Unlock()
	}
	return &schedule
}

func (app *App) StartScanner() {
	app.logger.Info("Starting scanner...")
	scanner := &Scanner{
		interval: app.config.ScanInterval,
		queue:    []time.Time{},
		logger:   app.logger,
	}

	scanner.StartScheduler()
	
	go func() {
		for {
			schedule := scanner.Next()
			app.logger.Info("New task spawned at", zap.Time("time", *schedule))
			if scanner.lastExecution != nil {
				execTime := time.Since(*scanner.lastExecution)
				app.logger.Info("Time since last executed task", zap.Duration("time", time.Duration(execTime)))
			}

			scanner.mu.Lock()
			scanner.logger.Info("Remaining tasks: ", zap.Int("count", len(scanner.queue)))
			scanner.lastExecution = schedule
			scanner.mu.Unlock()

			wg := new(sync.WaitGroup)
			wg.Add(len(app.config.Storages))

			for _, storage := range app.config.Storages {
				go storage.Scan(app, wg)
			}

			wg.Wait()
		}
	}()

	// Waitgroup?
	// Spawn gorutines for each file system
}
