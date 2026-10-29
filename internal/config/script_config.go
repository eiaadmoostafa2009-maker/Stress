package config

import (
	"fmt"
	"net/url"
	"time"
)

type ScriptConfig struct {
	Name      string
	BaseURL   string
	Users     int
	Duration  time.Duration
	Timeout   time.Duration
	ThinkTime time.Duration
}

func (c ScriptConfig) ScriptValidate() error {
	if c.Name == "" {
		return fmt.Errorf("config.name is required")
	}

	parsedURL, err := url.ParseRequestURI(c.BaseURL)
	if err != nil || parsedURL.Host == "" {
		return fmt.Errorf("config.base_url must be an absolute URL")
	}

	if parsedURL.Scheme != "http" &&
		parsedURL.Scheme != "https" {
		return fmt.Errorf("config.base_url must use HTTP or HTTPS")
	}

	if c.Users < 1 || c.Users > 100 {
		return fmt.Errorf("config.users must be between 1 and 100")
	}

	if c.Duration <= 0 || c.Duration > 5*time.Minute {
		return fmt.Errorf("config.duration_seconds must be between 1 and 300")
	}

	if c.Timeout <= 0 || c.Timeout > 60*time.Second {
		return fmt.Errorf("config.timeout_seconds must be between 1 and 60")
	}

	if c.ThinkTime < 0 || c.ThinkTime > 5*time.Minute {
		return fmt.Errorf("config.think_time_seconds must be between 0 and 60")
	}

	return nil
}
