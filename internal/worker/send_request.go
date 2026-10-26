package worker

import (
	"context"
	"io"
	"net/http"
	"stress/internal/model"
	"strings"
	"time"
)

func(w *worker) SendRequest(ctx context.Context, client *http.Client) (model.RequestResult, error) {
	start := time.Now()

	//validate the config
	if err := w.cfg.Validate(); err != nil {
		return model.RequestResult{}, err
	}

	//start a request
	req, err := http.NewRequestWithContext(ctx, w.cfg.Method, w.cfg.TargetUrl, strings.NewReader(w.cfg.Body))

	if w.cfg.JWT != "" {
       req.Header.Set(
         "Authorization",
         "Bearer "+w.cfg.JWT,
        )
    }
 
    if w.cfg.Body != "" {
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
		}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return model.RequestResult{
			Failed:  true,
			Error:   err.Error(),
			Latency: time.Since(start),
		}, err
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

	return requestResult, nil
}
