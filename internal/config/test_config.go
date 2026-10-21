package config

import (
	"fmt"
	"net/url"
	"time"
)

type TestConfig struct {
	TargetUrl   string
	Concurrency uint
	Duration    time.Duration
	Timeout     time.Duration
	JWT string
	Method string
	Body string
}

func (c *TestConfig) Validate() error {
	u, err := url.ParseRequestURI(c.TargetUrl)
	if err != nil {
		return err
	}
	// Ensure it has a scheme (http/https) and a host
	if u.Scheme == "" || u.Host == "" {
		return err
	}

	// Optional: Restrict to web-only traffic
	if u.Scheme != "http" && u.Scheme != "https" {
		return err
	}

	if c.Concurrency == 0 || c.Concurrency > 500 {
		return fmt.Errorf("Concurrency must be greater than 1 and less than or equal to 500")
	}

	if c.Duration <= 0 || c.Duration > 5*time.Minute {
		return fmt.Errorf("Duration must be greater than 0 and less than or equal to 5 minutes")
	}

	return nil
}
