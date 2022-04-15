package net

import (
	"net"
)

type Addr struct {
	Network string `json:"network"`
	Value   string `json:"value"`
}

func (a *Addr) Equal(b *Addr) bool {
	return a.Network == b.Network && a.Value == b.Value
}

type Interface struct {
	// Positive integer that starts at one, zero is never used
	Index     int    `json:"index"`
	MTU       int    `json:"mtu"`
	Name      string `json:"name"`
	Addresses []Addr `json:"addresses"`
}

func (i *Interface) Equal(j *Interface) bool {
	if len(i.Addresses) != len(j.Addresses) {
		return false
	}
	for i, lhs := range i.Addresses {
		rhs := &j.Addresses[i]
		if !lhs.Equal(rhs) {
			return false
		}
	}
	return i.Index == j.Index && i.MTU == j.MTU && i.Name == j.Name
}

func adaptAddrs(addrs []net.Addr) []Addr {
	result := make([]Addr, 0, len(addrs))
	for _, addr := range addrs {
		result = append(result, Addr{
			Network: addr.Network(),
			Value:   addr.String(),
		})
	}
	return result
}

func adaptInterfaces(netInterfaces []net.Interface, err error) ([]Interface, error) {
	result := make([]Interface, 0, len(netInterfaces))
	for _, netInterface := range netInterfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			return nil, err
		}
		result = append(result, Interface{
			Index:     netInterface.Index,
			MTU:       netInterface.MTU,
			Name:      netInterface.Name,
			Addresses: adaptAddrs(addrs),
		})
	}
	return result, nil
}

func Interfaces() ([]Interface, error) {
	return adaptInterfaces(net.Interfaces())
}
