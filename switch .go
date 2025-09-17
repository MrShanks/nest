package main

import (
	"fmt"
	"net"
	"sync"
)

type Switch struct {
	Name            string
	mu              sync.RWMutex
	ports           map[MAC]chan *Packet
	allInboundChans []chan *Packet
}

func NewSwitch(name string) *Switch {
	fmt.Printf("Configuring new Switch [%s]\n", name)
	return &Switch{
		Name:            name,
		ports:           make(map[MAC]chan *Packet),
		allInboundChans: make([]chan *Packet, 0),
	}
}

func (s *Switch) Plug(nic *NIC) {
	fmt.Printf("Connecting [%s] to [%s]\n", s.Name, net.HardwareAddr(nic.Mac[:]).String())

	s.mu.Lock()
	s.allInboundChans = append(s.allInboundChans, nic.InBound)
	s.ports[nic.Mac] = nic.InBound
	s.mu.Unlock()

	go func() {
		for frame := range nic.OutBound {
			s.forward(frame)
		}
	}()
}

func (s *Switch) forward(packet *Packet) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	destMac := packet.Header.DstMAC
	broadcastMAC := MAC{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	if destMac == broadcastMAC {
		fmt.Printf("[%s] Received broadcast. Floodig frame.\n", s.Name)
		for _, port := range s.allInboundChans {
			port <- packet
		}
		return
	}

	if destPort, ok := s.ports[destMac]; ok {
		fmt.Printf("[%s] Found destination [%s]. Forwarding fram.\n", s.Name, net.HardwareAddr(destMac[:]).String())
		destPort <- packet
	} else {
		fmt.Printf("[%s] Couldn't find destination mac address [%s]\n", s.Name, net.HardwareAddr(destMac[:]).String())
	}
}
