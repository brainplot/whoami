package api_test

import (
	"net/url"
	"testing"

	"github.com/desotech-it/whoami/api"
	"github.com/desotech-it/whoami/status"
)

func valuesFromStatusString(status string) url.Values {
	v := url.Values{}
	v.Set(api.HTTPParamStatus, status)
	return v
}

func TestParseStatusInValues(t *testing.T) {
	testCases := []struct {
		name   string
		values url.Values
		status status.Status
		err    error
	}{
		{
			name:   "Empty",
			values: url.Values{},
			status: "",
			err:    status.ErrInvalid,
		},
		{
			name:   "Up",
			values: valuesFromStatusString("up"),
			status: status.Up,
			err:    nil,
		},
		{
			name:   "Down",
			values: valuesFromStatusString("down"),
			status: status.Down,
			err:    nil,
		},
		{
			name:   "Invalid",
			values: valuesFromStatusString("invalid"),
			status: "",
			err:    status.ErrInvalid,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			status, err := api.ParseStatusInValues(tC.values)
			if got, want := status, tC.status; got != want {
				t.Errorf("got = %v; want = %v", got, want)
			}
			if got, want := err, tC.err; got != want {
				t.Errorf("got = %v; want = %v", got, want)
			}
		})
	}
}
