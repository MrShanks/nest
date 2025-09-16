package main

import (
	"fmt"
	"sync"
)

type Frame struct {
	SrcMAC  [6]byte
	DstMAC  [6]byte
	Payload []byte
}

type Switch struct {
	Name  string
	mu    sync.RWMutex
	ports map[[6]byte]chan *Frame
}

func NewSwitch(name string) *Switch {
	fmt.Printf("Configuring new Switch [%s]\n", name)
	return &Switch{
		Name:  name,
		ports: make(map[[6]byte]chan *Frame),
	}
}

func (s *Switch) Plug(nic *NIC) {
	s.ports[nic.MAC] = nic.InBound
	fmt.Printf("Connecting [%s] to [%v]\n", s.Name, nic.MAC)
	go func() {
		for frame := range nic.OutBound {
			s.forward(frame)
		}
	}()
}

func (s *Switch) forward(frame *Frame) {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, ok := s.ports[frame.DstMAC]; !ok {
		fmt.Printf("couldn't find destination mac address [%v]\n", frame.DstMAC)
		return
	}
	s.ports[frame.DstMAC] <- frame
}

type Host struct {
	Name string
	Nics map[[6]byte]*NIC
}

func NewHost(name string) *Host {
	fmt.Printf("Configuring new Host [%s]\n", name)
	return &Host{
		Name: name,
		Nics: make(map[[6]byte]*NIC),
	}
}

func (h *Host) ConfigureNIC(mac [6]byte) {
	fmt.Printf("Configuring NIC [%v] on [%s]\n", mac, h.Name)
	h.Nics[mac] = &NIC{
		MAC:      mac,
		OutBound: make(chan *Frame),
		InBound:  make(chan *Frame),
	}
}

func (h *Host) Send(frame *Frame) {
	fmt.Printf("[%s] sending a packet over interface [%v] to [%v]\n", h.Name, frame.SrcMAC, frame.DstMAC)
	h.Nics[frame.SrcMAC].OutBound <- frame
}

type NIC struct {
	MAC      [6]byte
	OutBound chan *Frame
	InBound  chan *Frame
}

func (n *NIC) Listen() {
	go func() {
		for frame := range n.InBound {
			fmt.Printf("--- NIC [%v] received frame: %+v ---\n", n.MAC, frame)
		}
	}()
}

func main() {
	s := NewSwitch("switch-1")

	A := NewHost("host-A")
	B := NewHost("host-B")

	A.ConfigureNIC([6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
	B.ConfigureNIC([6]byte{0x45, 0xB1, 0xC3, 0x01, 0xEE, 0xF2})

	A.Nics[[6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}].Listen()
	B.Nics[[6]byte{0x45, 0xB1, 0xC3, 0x01, 0xEE, 0xF2}].Listen()

	s.Plug(A.Nics[[6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}])
	s.Plug(B.Nics[[6]byte{0x45, 0xB1, 0xC3, 0x01, 0xEE, 0xF2}])

	frame := &Frame{
		SrcMAC:  [6]byte{0x45, 0xB1, 0xC3, 0x01, 0xEE, 0xF2},
		DstMAC:  [6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
		Payload: []byte("Hi this is a message for host A"),
	}

	B.Send(frame)

	select {}
}
