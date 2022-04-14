package memory

import (
	"context"

	"github.com/shirou/gopsutil/v3/mem"
)

type VirtualMemoryStat struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`
}

func adapt(stat *mem.VirtualMemoryStat, err error) (*VirtualMemoryStat, error) {
	if err != nil {
		return nil, err
	}
	result := &VirtualMemoryStat{
		Total:       stat.Total,
		Available:   stat.Available,
		Used:        stat.Used,
		UsedPercent: stat.UsedPercent,
	}
	return result, nil
}

func VirtualMemory() (*VirtualMemoryStat, error) {
	return adapt(mem.VirtualMemory())
}

func VirtualMemoryWithContext(ctx context.Context) (*VirtualMemoryStat, error) {
	return adapt(mem.VirtualMemoryWithContext(ctx))
}
