package engine

import (
	"context"
	"net/http"
	"stress/internal/config"
	"stress/internal/model"
	"stress/internal/worker"
)
type Engine interface{
    Load(parentctx context.Context) (model.TestResult, error)
}

type engine struct {
	cfg config.TestConfig
	client *http.Client
	worker worker.Worker
}

func NewEngine(cfg config.TestConfig, worker worker.Worker) Engine {
	return &engine{
		cfg: cfg,
		worker: worker,
	}
}

