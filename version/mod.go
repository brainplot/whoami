package version

import (
	"fmt"
	"strconv"
	"time"
)

var (
	major     string = "0"
	minor     string = "0"
	patch     string = "0"
	commit    string = "unknown"
	buildDate string = "1970-01-01T00:00:00Z"
)

type Info struct {
	Major     uint16
	Minor     uint16
	Patch     uint16
	Commit    string
	BuildDate time.Time
}

func parseVersionNumber(numStr string) (uint16, error) {
	result, err := strconv.ParseUint(numStr, 10, 16)
	return uint16(result), err
}

func MakeInfoFromStringComponents(major, minor, patch, commit, buildDate string) (*Info, error) {
	majorParsed, err := parseVersionNumber(major)
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", err.Error())
	}
	minorParsed, err := parseVersionNumber(minor)
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", err.Error())
	}
	patchParsed, err := parseVersionNumber(patch)
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", err.Error())
	}
	buildDateParsed, err := time.Parse(time.RFC3339, buildDate)
	if err != nil {
		return nil, fmt.Errorf("invalid build date: %s", err.Error())
	}
	info := &Info{
		Major:     majorParsed,
		Minor:     minorParsed,
		Patch:     patchParsed,
		Commit:    commit,
		BuildDate: buildDateParsed,
	}
	return info, nil
}

func ReadInfo() (*Info, error) {
	return MakeInfoFromStringComponents(major, minor, patch, commit, buildDate)
}
