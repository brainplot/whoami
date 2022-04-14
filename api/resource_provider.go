package api

import (
	"context"

	"github.com/desotech-it/whoami/api/memory"
)

type VirtualMemoryProvider interface {
	VirtualMemory(ctx context.Context) (*memory.VirtualMemoryStat, error)
}

type VirtualMemoryProviderFunc func(ctx context.Context) (*memory.VirtualMemoryStat, error)

func (f VirtualMemoryProviderFunc) VirtualMemory(ctx context.Context) (*memory.VirtualMemoryStat, error) {
	return f(ctx)
}
