package service

import (
	"context"
	"net/http"
	"time"
)

func(s *service) SendRequest(ctx context.Context) error{
	//config the client
    client := &http.Client{
		Timeout: 15*time.Second,
		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			return http.ErrUseLastResponse
		},
	}

	//use the worker
	resualt, err := s.worker.SendRequest(ctx, client)
	if err != nil {
		return err
	}
    
	//writing the report
	s.reporter.PrintReqReport(resualt)
	//return the resualts
	return nil
}