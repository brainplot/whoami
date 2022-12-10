package status_test

import (
	"testing"

	"github.com/desotech-it/whoami/status"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name   string
		result status.Status
		err    error
	}{
		{
			name:   "up",
			result: status.Up,
			err:    nil,
		},
		{
			name:   "down",
			result: status.Down,
			err:    nil,
		},
		{
			name:   "UP",
			result: status.Up,
			err:    nil,
		},
		{
			name:   "DOWN",
			result: status.Down,
			err:    nil,
		},
		{
			name:   "invalid",
			result: "",
			err:    status.ErrInvalid,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			status, err := status.Parse(tC.name)
			if got, want := status, tC.result; got != want {
				t.Errorf("got = %q; want = %q", got, want)
			}
			if got, want := err, tC.err; got != want {
				t.Errorf("got = %v; want = %v", got, want)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name   string
		status status.Status
	}{
		{
			name:   "up",
			status: status.Up,
		},
		{
			name:   "down",
			status: status.Down,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if got, want := tC.status.String(), tC.name; got != want {
				t.Errorf("got = %q; want = %q", got, want)
			}
		})
	}
}
