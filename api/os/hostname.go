package os

import "os"

type HostnameProvider interface {
	Hostname() (*HostnameInfo, error)
}

type HostnameProviderFunc func() (*HostnameInfo, error)

func (f HostnameProviderFunc) Hostname() (*HostnameInfo, error) {
	return f()
}

type HostnameInfo struct {
	Hostname string `json:"hostname"`
}

func Hostname() (*HostnameInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &HostnameInfo{Hostname: hostname}, nil
}
