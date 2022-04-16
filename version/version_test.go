package version_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/desotech-it/whoami/version"
)

var (
	commit = "feedcoffe"

	correctMajor     uint16    = 1
	correctMinor     uint16    = 2
	correctPatch     uint16    = 3
	correctBuildDate time.Time = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	wrongMajorStr     = "-1"
	wrongMinorStr     = "-2"
	wrongPatchStr     = "-3"
	wrongBuildDateStr = "test"
)

var (
	correctMajorStr     = formatVersionNumberComponent(correctMajor)
	correctMinorStr     = formatVersionNumberComponent(correctMinor)
	correctPatchStr     = formatVersionNumberComponent(correctPatch)
	correctBuildDateStr = formatBuildDate(correctBuildDate)
)

func formatVersionNumberComponent(value uint16) string {
	return strconv.FormatUint(uint64(value), 10)
}

func formatBuildDate(value time.Time) string {
	return value.Format(time.RFC3339)
}

func assertVersionNumberComponentEq(got uint16, want uint16, t *testing.T) {
	if got != want {
		t.Errorf("got = %d; want = %d", got, want)
	}
}

func TestMakeInfoFromStringComponents(t *testing.T) {
	t.Run("WrongMajor=error", func(t *testing.T) {
		info, err := version.MakeInfoFromStringComponents(
			wrongMajorStr,
			correctMinorStr,
			correctPatchStr,
			commit,
			correctBuildDateStr,
		)
		if err == nil {
			t.Error("got = nil; want = error")
		}
		if info != nil {
			t.Errorf("got = %v; want = nil", info)
		}
	})

	t.Run("WrongMinor=error", func(t *testing.T) {
		info, err := version.MakeInfoFromStringComponents(
			correctMajorStr,
			wrongMinorStr,
			correctPatchStr,
			commit,
			correctBuildDateStr,
		)
		if err == nil {
			t.Error("got = nil; want = error")
		}
		if info != nil {
			t.Errorf("got = %v; want = nil", info)
		}
	})

	t.Run("WrongPatch=error", func(t *testing.T) {
		info, err := version.MakeInfoFromStringComponents(
			correctMajorStr,
			correctMinorStr,
			wrongPatchStr,
			commit,
			correctBuildDateStr,
		)
		if err == nil {
			t.Error("got = nil; want = error")
		}
		if info != nil {
			t.Errorf("got = %v; want = nil", info)
		}
	})

	t.Run("WrongBuildDate=error", func(t *testing.T) {
		info, err := version.MakeInfoFromStringComponents(
			correctMajorStr,
			correctMinorStr,
			correctPatchStr,
			commit,
			wrongBuildDateStr,
		)
		if err == nil {
			t.Error("got = nil; want = error")
		}
		if info != nil {
			t.Errorf("got = %v; want = nil", info)
		}
	})

	t.Run("CorrectInfo=OK", func(t *testing.T) {
		info, err := version.MakeInfoFromStringComponents(
			correctMajorStr,
			correctMinorStr,
			correctPatchStr,
			commit,
			correctBuildDateStr,
		)
		if err != nil {
			t.Errorf("got = %v; want = nil", err)
		}
		assertVersionNumberComponentEq(info.Major, correctMajor, t)
		assertVersionNumberComponentEq(info.Minor, correctMinor, t)
		assertVersionNumberComponentEq(info.Patch, correctPatch, t)
		if !info.BuildDate.Equal(correctBuildDate) {
			t.Errorf("got = %v; want = %v", info.BuildDate, correctBuildDate)
		}
	})
}
