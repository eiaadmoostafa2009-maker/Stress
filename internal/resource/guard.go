package resource

import (
	"context"
	"fmt"
)

type GuardConfig struct {
	MaxCPUPercent float64
	MaxRAMPercent float64
}

func (g GuardConfig) Check(snapshot Snapshot) error {
	if snapshot.CPUPercent >= g.MaxCPUPercent {
		return fmt.Errorf(
			"CPU usage reached %.1f%%",
			snapshot.CPUPercent,
		)
	}

	if snapshot.SystemRAMUsed >= g.MaxRAMPercent {
		return fmt.Errorf(
			"system RAM usage reached %.1f%%",
			snapshot.SystemRAMUsed,
		)
	}

	return nil
}

func WatchResources(ctx context.Context, samples <-chan Snapshot, cfg GuardConfig, cancel context.CancelFunc) <-chan error {
	stopped := make(chan error, 1)

	go func() {
		defer close(stopped)

		for {
			select {
			case <-ctx.Done():
				return

			case snapshot, ok := <-samples:
				if !ok {
					return
				}

				if err := cfg.Check(snapshot); err != nil {
					stopped <- err
					cancel()
					return
				}
			}
		}
	}()

	return stopped
}
