package loadtest

import (
	"context"
	"io"
	"net/http"
	"stress/internal/config"
	"stress/internal/model"
	"strings"
	"time"
)

func SendRequest(ctx context.Context, cfg *config.TestConfig, client *http.Client) model.RequestResult {
	start := time.Now()

	//validate the config
	if err := cfg.Validate(); err != nil {
		return model.RequestResult{}
	}

	//start a request
	req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.TargetUrl, strings.NewReader(cfg.Body))

	if cfg.JWT != "" {
       req.Header.Set(
         "Authorization",
         "Bearer "+cfg.JWT,
        )
    }
 
    if cfg.Body != "" {
       req.Header.Set(
         "Content-Type",
         "application/json",
        )
    }

	if err != nil {
		return model.RequestResult{
			Failed:  true,
			Error:   err.Error(),
			Latency: time.Since(start),
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return model.RequestResult{
			Failed:  true,
			Error:   err.Error(),
			Latency: time.Since(start),
		}
	}
	defer resp.Body.Close()

	_, _ = io.Copy(
		io.Discard,
		io.LimitReader(resp.Body, 64*1024),
	)

	latency := time.Since(start)

	//create a request result
	requestResult := model.RequestResult{
		StatusCode: resp.StatusCode,
		Latency:    latency,
		Failed:     resp.StatusCode >= 500,
	}

	return requestResult
}
