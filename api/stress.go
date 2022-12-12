package api

import (
	"context"
	"time"

	"github.com/desotech-it/whoami/api/memory"
)

type Stresser interface {
	Stress(context.Context) error
}

type StressSession struct {
	StartedAt  time.Time  `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
}

type MemoryStressSession struct {
	StressSession
	BytesAllocated int `json:"bytesAllocated"`
	stresser       *memory.Stresser
	cancelFunc     context.CancelFunc
}
