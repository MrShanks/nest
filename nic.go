package main

import "fmt"

type NetworkInterface struct {
	MACaddr  string
	inbound  chan *Frame
	outbound chan *Frame
}

func NewNetworkInterface(mac string) *NetworkInterface {
	return &NetworkInterface{
		MACaddr: mac,
		inbound: make(chan *Frame, 64),
	}
}

func (nic *NetworkInterface) Connect(switchPort chan *Frame) {
	fmt.Printf("Connecting to Switch port\n")
	nic.outbound = switchPort
}

func (nic *NetworkInterface) Send(frame *Frame) {
	fmt.Printf("Sending frame trough nic outbound channel\n")
	nic.outbound <- frame
}

func (nic *NetworkInterface) Receive() <-chan *Frame {
	fmt.Printf("Receiving from switch towards inbound nic channel\n")
	return nic.inbound
}
