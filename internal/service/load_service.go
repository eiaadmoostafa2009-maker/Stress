package service

import (
	"context"
	"fmt"
)

func (s *service) Load(ctx context.Context) error {
	if err := s.cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	result, err := s.engine.Load(ctx)
	if err != nil {
		return fmt.Errorf("load test failed: %w", err)
	}

	s.reporter.PrintLoadReport(result)
	return nil
}
