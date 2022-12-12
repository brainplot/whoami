package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/desotech-it/whoami/api/memory"
)

const (
	HTTPParamInterval  = "interval"
	HTTPParamAllocSize = "allocSize"
)

type Stresser interface {
	Stress(context.Context) error
}

type stressSession struct {
	StartedAt  time.Time
	FinishedAt *time.Time
	cancelFunc context.CancelFunc
}

type memoryStressSession struct {
	stressSession
	stresser *memory.Stresser
}

type memoryStressSessionResponse struct {
	StartedAt      time.Time  `json:"startedAt"`
	FinishedAt     *time.Time `json:"finishedAt,omitempty"`
	BytesAllocated int        `json:"bytesAllocated"`
}

func buildMemoryStressSessionResponse(s *memoryStressSession) memoryStressSessionResponse {
	return memoryStressSessionResponse{
		StartedAt:      s.StartedAt,
		FinishedAt:     s.FinishedAt,
		BytesAllocated: s.stresser.BytesAllocated(),
	}
}

func ParseMemoryStressParams(r *http.Request) (*memory.StressParameters, error) {
	requestInterval := r.FormValue(HTTPParamInterval)
	var interval time.Duration
	var allocSize int
	if len(requestInterval) != 0 {
		if parsed, err := time.ParseDuration(requestInterval); err != nil {
			return nil, err
		} else {
			interval = parsed
		}
	}
	requestAllocSize := r.FormValue(HTTPParamAllocSize)
	if len(requestAllocSize) != 0 {
		if parsed, err := strconv.ParseInt(requestAllocSize, 0, 0); err != nil {
			return nil, err
		} else {
			allocSize = int(parsed)
		}
	}
	return &memory.StressParameters{Interval: interval, AllocationSize: allocSize}, nil
}
