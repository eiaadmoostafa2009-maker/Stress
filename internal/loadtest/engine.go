package loadtest

import (
	"context"
	"fmt"
	"net/http"
	"stress/internal/config"
	"stress/internal/model"
	"stress/internal/resource"
	"sync"
	"time"
)

type Engine struct {
	Client *http.Client
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(parentctx context.Context, cfg config.TestConfig) (model.TestResult, error) {
	//validate config
	err := cfg.Validate()
	if err != nil {
		return model.TestResult{}, err
	}

	//config client
	client := &http.Client{
		Timeout: cfg.Timeout,
		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			return http.ErrUseLastResponse
		},
	}

	e.Client = client

	//define ctx
	testCtx, cancel := context.WithTimeout(parentctx, cfg.Duration)
	defer cancel()

	//start monitoring
	samples := resource.StartMonitor(testCtx)

	guardStopped := resource.WatchResources(testCtx, samples, resource.GuardConfig{
		MaxCPUPercent: 90,
		MaxRAMPercent: 90,
	}, cancel)

	//start the test
	start := time.Now()
	results := make(chan model.RequestResult, cfg.Concurrency)
	var workers sync.WaitGroup

	for i := 0; i < int(cfg.Concurrency); i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for {
				if testCtx.Err() != nil {
					return
				}

				result := SendRequest(testCtx, &cfg, e.Client)
				results <- result
			}
		}()

	}

	go func() {
		workers.Wait()
		close(results)
	}()

	var summary model.TestResult
	summary.Url = cfg.TargetUrl
	summary.StartedAt = start

	var totalLatency time.Duration

	// One goroutine owns all metric updates.
	for result := range results {
		summary.TotalRequests++

		if result.Failed {
			summary.FailureCount++
		}

		totalLatency += result.Latency
	}

	summary.Duration = time.Since(start)

	if summary.TotalRequests > 0 {
		summary.AverageLatency = int64(time.Duration(int64(totalLatency) / summary.TotalRequests))
	}

	if summary.Duration > 0 {
		summary.ReqPerSec = float64(summary.TotalRequests) / summary.Duration.Seconds()
	}

	if err := parentctx.Err(); err != nil {
		return summary, fmt.Errorf(
			"load test cancelled: %w",
			err,
		)
	}

	if guardStopped != nil {
		return summary, <-guardStopped
	}

	return summary, nil
}
