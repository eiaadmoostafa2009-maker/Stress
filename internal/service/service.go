package service

import (
	"context"

	"stress/internal/config"
	"stress/internal/engine"
	"stress/internal/reporter"
	"stress/internal/worker"
)

type Service interface{
    Load(ctx context.Context) error 
	SendRequest(ctx context.Context) error
}

type service struct {
	engine   engine.Engine
	reporter reporter.Reporter
	worker worker.Worker
	cfg config.TestConfig
}	

func NewService(engine engine.Engine, reporter reporter.Reporter, cfg config.TestConfig, worker worker.Worker) Service {
	
	return &service{
		engine:   engine,
		reporter: reporter,
		cfg: cfg,
		worker: worker,
	}
}
