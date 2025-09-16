package main

import (
	"net"
)

type NIC struct {
	Name       string
	MAC        [6]byte
	IP         [4]byte
	SocketPath string
	Socket     net.Listener
}

func NewNIC(name string, mac [6]byte, ip [4]byte) *NIC {
	return &NIC{
		Name: name,
		MAC:  mac,
		IP:   ip,
	}
}
