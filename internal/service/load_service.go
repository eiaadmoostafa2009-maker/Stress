package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"stress/internal/config"
	"stress/internal/loadtest"
)

type LoadTestService struct {
	engine   loadtest.Engine
	reporter io.Writer
}

func NewLoadTestService(engine loadtest.Engine, reporter io.Writer) *LoadTestService {
	if reporter == nil {
		reporter = os.Stdout
	}

	return &LoadTestService{
		engine:   engine,
		reporter: reporter,
	}
}

func (s *LoadTestService) Run(ctx context.Context, cfg config.TestConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	result, err := s.engine.Run(ctx, cfg)
	if err != nil {
		return fmt.Errorf("load test failed: %w", err)
	}

	loadtest.PrintReport(s.reporter, result)
	return nil
}
