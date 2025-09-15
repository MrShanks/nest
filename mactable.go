package main

import (
	"fmt"
	"sync"
)

type MACAddressTable struct {
	mu    sync.Mutex
	table map[string]chan *Frame
}

func NewMACAddressTable() *MACAddressTable {
	return &MACAddressTable{
		table: make(map[string]chan *Frame),
	}
}

func (t *MACAddressTable) Learn(mac string, port chan *Frame) {
	fmt.Printf("Learning a new MAC Address: %s\n", mac)
	t.mu.Lock()
	defer t.mu.Unlock()
	t.table[mac] = port
}

func (t *MACAddressTable) GetPort(mac string) (chan *Frame, bool) {
	fmt.Printf("Retrieving port for MAC Address: %s\n", mac)
	t.mu.Lock()
	defer t.mu.Unlock()
	port, found := t.table[mac]
	return port, found
}
