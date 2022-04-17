package version_test

import (
	"testing"
	"time"

	"github.com/desotech-it/whoami/version"
)

var (
	testVersion = version.Info{
		Major:     1,
		Minor:     2,
		Patch:     3,
		Commit:    "feedcoffe",
		BuildDate: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
)

func assertVersionNumberComponentEq(got uint16, want uint16, t *testing.T) {
	if got != want {
		t.Errorf("got = %d; want = %d", got, want)
	}
}

func TestMakeInfoFromStringComponents(t *testing.T) {
	testCases := []struct {
		name      string
		major     string
		minor     string
		patch     string
		commit    string
		buildDate string
		want      *version.Info
	}{
		{
			name:      "WrongMajor",
			major:     "-1",
			minor:     "2",
			patch:     "3",
			commit:    "feedcoffee",
			buildDate: "1970-01-01T00:00:00Z",
		},
		{
			name:      "WrongMinor",
			major:     "1",
			minor:     "-2",
			patch:     "3",
			commit:    "feedcoffee",
			buildDate: "1970-01-01T00:00:00Z",
		},
		{
			name:      "WrongPatch",
			major:     "1",
			minor:     "2",
			patch:     "-3",
			commit:    "feedcoffee",
			buildDate: "1970-01-01T00:00:00Z",
		},
		{
			name:      "WrongBuildDate",
			major:     "1",
			minor:     "2",
			patch:     "3",
			commit:    "feedcoffee",
			buildDate: "test",
		},
		{
			name:      "OK",
			major:     "1",
			minor:     "2",
			patch:     "3",
			commit:    "feedcoffe",
			buildDate: "1970-01-01T00:00:00Z",
			want:      &testVersion,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got, err := version.MakeInfoFromStringComponents(
				tC.major, tC.minor, tC.patch, tC.commit, tC.buildDate,
			)
			if got != nil {
				if err != nil {
					t.Error(err)
				}
				assertVersionNumberComponentEq(got.Major, tC.want.Major, t)
				assertVersionNumberComponentEq(got.Minor, tC.want.Minor, t)
				assertVersionNumberComponentEq(got.Patch, tC.want.Patch, t)
				if got.Commit != tC.want.Commit {
					t.Errorf("got = %v; want = %v", got.Commit, tC.want.Commit)
				}
				if !got.BuildDate.Equal(tC.want.BuildDate) {
					t.Errorf("got = %v; want = %v", got.BuildDate, tC.want.BuildDate)
				}
			} else {
				if err == nil {
					t.Error("both info and err are nil")
				}
			}
		})
	}
}
