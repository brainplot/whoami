package net_test

import (
	"testing"

	"github.com/desotech-it/whoami/api/net"
)

func TestAddrEqual(t *testing.T) {
	testCases := []struct {
		name string
		lhs  net.Addr
		rhs  net.Addr
		want bool
	}{
		{
			name: "Equal",
			lhs:  net.Addr{"test", "10.0.0.1"},
			rhs:  net.Addr{"test", "10.0.0.1"},
			want: true,
		},
		{
			name: "NameNotEqual",
			lhs:  net.Addr{"test1", "10.0.0.1"},
			rhs:  net.Addr{"test2", "10.0.0.1"},
			want: false,
		},
		{
			name: "ValueNotEqual",
			lhs:  net.Addr{"test", "10.0.0.1"},
			rhs:  net.Addr{"test", "10.0.0.2"},
			want: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := tC.lhs.Equal(&tC.rhs)
			want := tC.want
			if got != want {
				t.Errorf("got = %t; want = %t", got, want)
			}
		})
	}
}

func TestInterfaceEqual(t *testing.T) {
	testCases := []struct {
		name string
		lhs  net.Interface
		rhs  net.Interface
		want bool
	}{
		{
			name: "Equal",
			lhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			rhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			want: true,
		},
		{
			name: "IndexNotEqual",
			lhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			rhs: net.Interface{
				Index:     43,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			want: false,
		},
		{
			name: "MTUNotEqual",
			lhs: net.Interface{
				Index:     42,
				MTU:       4243,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			rhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			want: false,
		},
		{
			name: "NameNotEqual",
			lhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			rhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "testtest",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			want: false,
		},
		{
			name: "AddressesNotEqual",
			lhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			rhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "172.31.0.1"}, {"test2", "192.168.0.1"}},
			},
			want: false,
		},
		{
			name: "AddressesNotEqualSize",
			lhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}, {"test2", "172.31.0.1"}},
			},
			rhs: net.Interface{
				Index:     42,
				MTU:       4242,
				Name:      "test",
				Addresses: []net.Addr{{"test1", "10.0.0.1"}},
			},
			want: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := tC.lhs.Equal(&tC.rhs)
			want := tC.want
			if got != want {
				t.Errorf("got = %t; want = %t", got, want)
			}
		})
	}
}
