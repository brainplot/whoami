package memory

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"time"
)

const (
	DefaultStressInterval       = 1 * time.Second
	DefaultStressAllocationSize = 32 * 1024 * 1024
)

var (
	ErrNegativeStressInterval = errors.New("negative interval")
	ErrNegativeAllocationSize = errors.New("negative allocation size")
)

type StressParameters struct {
	Interval       time.Duration `json:"interval" form:"interval"`
	AllocationSize int           `json:"allocationSize" form:"allocationSize"`
}

type Stresser struct {
	params         StressParameters
	bytesAllocated int
	byteSource     io.Reader
	buff           [][]byte
}

func NewStresser(params StressParameters) (*Stresser, error) {
	byteSource := rand.New(rand.NewSource(time.Now().Unix()))
	return NewStresserWithByteSource(params, byteSource)
}

func NewStresserWithByteSource(params StressParameters, src io.Reader) (*Stresser, error) {
	if params.Interval < 0 {
		return nil, ErrNegativeStressInterval
	}
	if params.AllocationSize < 0 {
		return nil, ErrNegativeAllocationSize
	}
	if params.Interval == 0 {
		params.Interval = DefaultStressInterval
	}
	if params.AllocationSize == 0 {
		params.AllocationSize = DefaultStressAllocationSize
	}
	session := &Stresser{
		params:         params,
		bytesAllocated: 0,
		byteSource:     src,
		buff:           make([][]byte, 0),
	}
	return session, nil
}

func (s *Stresser) Stress(ctx context.Context) error {
	var err error = nil
	done := ctx.Done()
	select {
	case <-done:
		break
	default:
		ticker := time.NewTicker(s.params.Interval)
	loop:
		for range ticker.C {
			select {
			case <-done:
				break loop
			default:
				err = s.fillBuffer()
				if err != nil {
					break loop
				}
			}
		}
		ticker.Stop()
	}
	return err
}

func (s *Stresser) fillBuffer() error {
	chunk := make([]byte, s.params.AllocationSize)
	if n, err := s.byteSource.Read(chunk); err != nil {
		return err
	} else {
		s.buff = append(s.buff, chunk)
		s.bytesAllocated += n
		return nil
	}
}

func (s *Stresser) BytesAllocated() int {
	return s.bytesAllocated
}
