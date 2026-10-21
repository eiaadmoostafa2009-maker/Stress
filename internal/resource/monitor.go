package resource

import (
	"context"

	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Snapshot struct {
	CPUPercent     float64
	SystemRAMUsed  float64
	SystemRAMTotal uint64
	AvailableRAM   uint64
	ProcessRAM     uint
	GoHeapAlloc    uint64
}

func ReadSnapshot() (Snapshot, error) {
	cpuValues, err := cpu.Percent(0, false)
	if err != nil {
		return Snapshot{}, err
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return Snapshot{}, err
	}

	var runtimeMemory runtime.MemStats
	runtime.ReadMemStats(&runtimeMemory)

	var cpuPercent float64
	if len(cpuValues) > 0 {
		cpuPercent = cpuValues[0]
	}

	return Snapshot{
		CPUPercent:     cpuPercent,
		SystemRAMUsed:  memory.UsedPercent,
		AvailableRAM:   memory.Available,
		SystemRAMTotal: memory.Total,
		GoHeapAlloc:    runtimeMemory.HeapAlloc,
	}, nil
}

func StartMonitor(ctx context.Context) <-chan Snapshot {
	samples := make(chan Snapshot)

	go func() {
		defer close(samples)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				snapshot, err := ReadSnapshot()
				if err != nil {
					//guard
					continue
				}

				select {
				case samples <- snapshot:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return samples
}
