package memory

import "context"

type VirtualMemoryProvider interface {
	VirtualMemory(ctx context.Context) (*VirtualMemoryStat, error)
}

type VirtualMemoryProviderFunc func(ctx context.Context) (*VirtualMemoryStat, error)

func (f VirtualMemoryProviderFunc) VirtualMemory(ctx context.Context) (*VirtualMemoryStat, error) {
	return f(ctx)
}
