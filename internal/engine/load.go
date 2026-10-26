package engine

import (
	"context"
	"fmt"
	"net/http"
	"stress/internal/model"
	"stress/internal/resource"
	"sync"
	"time"
)

func (e *engine) Load(parentctx context.Context) (model.TestResult, error) {
	//validate config
	err := e.cfg.Validate()
	if err != nil {
		return model.TestResult{}, err
	}

	//config client
	client := &http.Client{
		Timeout: e.cfg.Timeout,
		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			return http.ErrUseLastResponse
		},
	}

	e.client = client

	//define ctx
	testCtx, cancel := context.WithTimeout(parentctx, e.cfg.Duration)
	defer cancel()

	//start monitoring
	samples := resource.StartMonitor(testCtx)

	guardStopped := resource.WatchResources(testCtx, samples, resource.GuardConfig{
		MaxCPUPercent: 90,
		MaxRAMPercent: 90,
	}, cancel)

	//start the test
	start := time.Now()
	results := make(chan model.RequestResult, e.cfg.Concurrency)
	var workers sync.WaitGroup

	for i := 0; i < int(e.cfg.Concurrency); i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for {
				if testCtx.Err() != nil {
					return
				}

				result, err := e.worker.SendRequest(testCtx, e.client)
				if err != nil {
		            return 
	            }
				results <- result
			}
		}()

	}

	go func() {
		workers.Wait()
		close(results)
	}()

	var summary model.TestResult
	summary.Url = e.cfg.TargetUrl
	summary.StartedAt = start

	var totalLatency time.Duration

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