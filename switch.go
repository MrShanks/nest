package main

import (
	"fmt"
	"net"
	"time"
)

type Switch struct {
	Name  string
	Ports []net.Listener
}

func NewSwitch(name string) *Switch {
	return &Switch{
		Name:  name,
		Ports: make([]net.Listener, 0),
	}
}

func (s *Switch) Boot() {
	fmt.Printf("[%s] starting up...\n", s.Name)
	for {
		time.Sleep(10 * time.Second)
	}
}
