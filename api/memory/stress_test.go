package memory_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/desotech-it/whoami/api/memory"
)

func TestStresserErr(t *testing.T) {
	testCases := []struct {
		name   string
		params memory.StressParameters
		err    error
	}{
		{
			name: "NegativeStressInterval",
			params: memory.StressParameters{
				Interval:       -1 * time.Second,
				AllocationSize: 42,
			},
			err: memory.ErrNegativeStressInterval,
		},
		{
			name: "NegativeAllocationSize",
			params: memory.StressParameters{
				Interval:       42 * time.Second,
				AllocationSize: -42,
			},
			err: memory.ErrNegativeAllocationSize,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if _, err := memory.NewStresser(tC.params); err != tC.err {
				t.Errorf("got = %v; want = %v", err, tC.err)
			}
		})
	}
}

func TestStresser(t *testing.T) {
	const (
		testString   = "test"
		testInterval = 25 * time.Millisecond
	)

	testParams := memory.StressParameters{Interval: testInterval, AllocationSize: 1}

	testCases := []struct {
		name           string
		duration       time.Duration
		reader         io.Reader
		bytesAllocated int
		stressError    error
	}{
		{
			name:           "EntireBuffer",
			duration:       time.Duration(len(testString)) * testInterval,
			reader:         strings.NewReader(testString),
			bytesAllocated: len(testString),
		},
		{
			name:           "HalfBuffer",
			duration:       time.Duration(len(testString)/2) * testInterval,
			reader:         strings.NewReader(testString),
			bytesAllocated: len(testString) / 2,
		},
		{
			name:           "EmptyReader",
			duration:       testInterval * 2,
			reader:         new(bytes.Buffer),
			bytesAllocated: 0,
			stressError:    io.EOF,
		},
		{
			name:           "ZeroDurationContext",
			duration:       0,
			reader:         new(bytes.Buffer),
			bytesAllocated: 0,
		},
		{
			name:           "ShortContext",
			duration:       testInterval / 2,
			reader:         new(bytes.Buffer),
			bytesAllocated: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if session, err := memory.NewStresserWithByteSource(testParams, tC.reader); err != nil {
				t.Error(err)
			} else {
				ctx, cancelFunc := context.WithTimeout(context.Background(), tC.duration)
				defer cancelFunc()
				if err := session.Stress(ctx); err != tC.stressError {
					t.Errorf("got = %v; want = %v", err, tC.stressError)
				}
				got := session.BytesAllocated()
				want := tC.bytesAllocated
				if got != want {
					t.Errorf("got = %d; want = %d", got, want)
				}
			}
		})
	}
}
