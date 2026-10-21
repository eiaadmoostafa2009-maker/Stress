package model

import (
	"time"
)

type TestResult struct {
	Url                 string
	SuccessCount        int
	FailureCount        int
	StartedAt           time.Time
	Duration            time.Duration
	AverageLatency      int64
	TotalRequests       int64
	ReqPerSec           float64
	AverageResponseTime float64
	AverageSuccessRate  float64
}
