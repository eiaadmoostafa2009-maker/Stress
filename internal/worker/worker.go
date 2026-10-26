package worker
import(
    "stress/internal/config"
	"net/http"
	"context"
	"stress/internal/model"
)

type Worker interface{
SendRequest(ctx context.Context, client *http.Client) (model.RequestResult, error) 
}

type worker struct{
	cfg config.TestConfig
}

func NewWorker(cfg config.TestConfig) Worker{
	return &worker{
		cfg: cfg,
	}
}