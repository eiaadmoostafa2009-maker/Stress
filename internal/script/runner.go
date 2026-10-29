package script

import (
	"context"
	"fmt"
	"sync"
)

func Run(ctx context.Context, path string) error {
	cfg, err := LoadConfig(path)
	if err != nil {
		return err
	}

	testCtx, cancel := context.WithTimeout(ctx, cfg.Duration)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, cfg.Users)

	for i := 0; i < cfg.Users; i++ {
		wg.Add(1)
		go func() {
			err := RunVirtualUser(testCtx, cfg, path)
			if err != nil {
				errChan <- err
				cancel()
			}
		}()

	}
	wg.Wait()
	defer close(errChan)

	for err := range errChan {
		return fmt.Errorf("Error: %w", err)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
