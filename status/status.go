package status

import (
	"errors"
	"strings"
)

var (
	ErrInvalid = errors.New("invalid status")
)

type Status string

func Parse(status string) (Status, error) {
	status = strings.ToLower(status)
	if status == "up" {
		return Up, nil
	}
	if status == "down" {
		return Down, nil
	}
	return "", ErrInvalid
}

func (s Status) String() string {
	return string(s)
}

const (
	Up   Status = "up"
	Down Status = "down"
)

type InstanceStatus struct {
	Health    Status
	Readiness Status
}

var Current = InstanceStatus{Health: Up, Readiness: Up}
