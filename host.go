package main

import (
	"fmt"
	"sync"
)

type Host struct {
	Name         string
	Interfaces   []*NetworkInterface
	shutdownChan chan struct{}
	wg           sync.WaitGroup
}

func NewHost(name string, numInterfaces int) *Host {
	fmt.Printf("Creating new Host [%s]\n", name)
	h := &Host{
		Name:         name,
		Interfaces:   make([]*NetworkInterface, 0),
		shutdownChan: make(chan struct{}),
	}

	for range numInterfaces {
		mac := fmt.Sprintf("00:00:00:00:00:%02d", len(h.Interfaces)+1)
		nic := NewNetworkInterface(mac)
		h.Interfaces = append(h.Interfaces, nic)

		h.wg.Add(1)
		go h.listenOn(nic)
	}

	return h
}

func (h *Host) listenOn(nic *NetworkInterface) {
	defer h.wg.Done()
	fmt.Printf("[%s] Listening on interface %s\n", h.Name, nic.MACaddr)
	for {
		select {
		case frame := <-nic.inbound:
			fmt.Printf("[%s] Received frame for %s on interface %s. Payload: %s\n", h.Name, frame.DestinationMAC, nic.MACaddr, string(frame.Payload))
		case <-h.shutdownChan:
			fmt.Printf("[%s] Stopping listener on interface %s\n", h.Name, nic.MACaddr)
			return
		}
	}
}

func (h *Host) SendData(interfaceIndex int, destMAC string, data []byte) error {
	if interfaceIndex < 0 || interfaceIndex >= len(h.Interfaces) {
		return fmt.Errorf("invalid interface index")
	}

	nic := h.Interfaces[interfaceIndex]
	frame := &Frame{
		SourceMAC:      nic.MACaddr,
		DestinationMAC: destMAC,
		Payload:        data,
	}

	nic.Send(frame)
	return nil
}

func (h *Host) Stop() {
	close(h.shutdownChan)
	h.wg.Wait()
}

func (h *Host) AddNetworkInterface(mac string) {
	nic := NewNetworkInterface(mac)
	h.Interfaces = append(h.Interfaces, nic)
}
