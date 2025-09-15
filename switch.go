package main

import (
	"fmt"
	"sync"
	"time"
)

type Switch struct {
	Name            string
	ports           []chan *Frame
	macAddressTable *MACAddressTable
	shutdownChan    chan struct{}
	wg              sync.WaitGroup
}

func NewSwitch(name string) *Switch {
	fmt.Printf("Creating new Switch [%s]\n", name)
	s := &Switch{
		Name:            name,
		ports:           make([]chan *Frame, 0),
		macAddressTable: NewMACAddressTable(),
		shutdownChan:    make(chan struct{}),
	}

	s.wg.Add(1)
	go s.start()

	return s
}

func (s *Switch) ConfigureNewPort() chan *Frame {
	port := make(chan *Frame, 64)
	fmt.Printf("Configuring new port [%v]\n", port)
	s.ports = append(s.ports, port)
	return port
}

func (s *Switch) start() {
	defer s.wg.Done()
	fmt.Printf("[%s] Switch started\n", s.Name)

	for {
		select {
		case <-s.shutdownChan:
			fmt.Printf("[%s] Switch powering off\n", s.Name)
			return
		default:
			for i, port := range s.ports {
				select {
				case frame := <-port:
					fmt.Printf("[%s] Received frame from %s to %s on port %d\n", s.Name, frame.SourceMAC, frame.DestinationMAC, i)

					s.macAddressTable.Learn(frame.SourceMAC, port)
					s.forward(frame, port)
				default:
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (s *Switch) forward(frame *Frame, incomingPort chan *Frame) {
	destPort, found := s.macAddressTable.GetPort(frame.DestinationMAC)

	if found {
		if destPort != incomingPort {
			fmt.Printf("[%s] Forwarding frame from %s to %s via known port\n", s.Name, frame.SourceMAC, frame.DestinationMAC)
			destPort <- frame
		}
	} else {
		fmt.Printf("[%s] Flooding frame from %s to %s\n", s.Name, frame.SourceMAC, frame.DestinationMAC)
		for _, p := range s.ports {
			if p != incomingPort {
				p <- frame
			}
		}
	}
}

func (s *Switch) Stop() {
	close(s.shutdownChan)
	s.wg.Wait()
}
