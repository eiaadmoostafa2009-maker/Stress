package model

import "time"

type RequestResult struct {
	StatusCode int
	Latency    time.Duration
	Failed     bool
	Error      string
}
