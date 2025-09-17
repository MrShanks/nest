package main

import (
	"fmt"
	"net"
	"sync"
)

type Host struct {
	Name        string
	Nics        map[MAC]*NIC
	mu          sync.Mutex
	ArpCache    map[IP]MAC
	packetQueue map[IP]*Packet
}

func NewHost(name string) *Host {
	fmt.Printf("Configuring new Host [%s]\n", name)
	return &Host{
		Name:        name,
		Nics:        make(map[MAC]*NIC),
		ArpCache:    make(map[IP]MAC),
		packetQueue: make(map[IP]*Packet),
	}
}

func (h *Host) ConfigureNIC(mac MAC, ip IP) {
	fmt.Printf("Configuring NIC [%s] on [%s]\n", net.HardwareAddr(mac[:]).String(), h.Name)
	h.Nics[mac] = &NIC{
		Mac:      mac,
		Ip:       ip,
		OutBound: make(chan *Packet),
		InBound:  make(chan *Packet),
		Host:     h,
	}
}

func (h *Host) Send(data []byte, dstIP IP, srcMAC MAC) {
	fmt.Printf("[%s] sending a packet over interface [%s] to [%s]\n", h.Name, net.HardwareAddr(srcMAC[:]).String(), net.IP(dstIP[:]).String())

	packet := &Packet{
		Header: EthernetHeader{
			SrcMAC: srcMAC,
		},
		Payload: IPv4Packet{
			DstIP: dstIP,
			SrcIP: h.Nics[srcMAC].Ip,
			Data:  data,
		},
	}

	h.mu.Lock()
	dstMAC, ok := h.ArpCache[dstIP]
	h.mu.Unlock()
	if ok {
		// If we know the MAC, send it directly
		fmt.Printf("[%s] Destination MAC found in cache. Sending packet.\n", h.Name)
		packet.Header.DstMAC = dstMAC
		h.Nics[srcMAC].OutBound <- packet
	} else {
		// If we don't know the MAC, queue the packet and send an ARP request
		fmt.Printf("[%s] Destination MAC not in cache. Queuing packet and sending ARP request.\n", h.Name)
		h.mu.Lock()
		h.packetQueue[dstIP] = packet
		h.mu.Unlock()

		arpRequest := &Packet{
			Header: EthernetHeader{
				SrcMAC: srcMAC,
				DstMAC: MAC{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // Broadcast
			},
			Payload: IPv4Packet{
				DstIP: dstIP,
				SrcIP: h.Nics[srcMAC].Ip,
			},
		}
		h.Nics[srcMAC].OutBound <- arpRequest
	}
}

func (h *Host) updateArpCache(ip IP, mac MAC) {
	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Printf("[%s] Updating ARP cache: %s -> %s\n", h.Name, net.IP(ip[:]).String(), net.HardwareAddr(mac[:]).String())
	h.ArpCache[ip] = mac

	// After updating cache, check if any packets were waiting and send them
	if queuedPacket, ok := h.packetQueue[ip]; ok {
		fmt.Printf("[%s] Found queued packet for %s. Sending now.\n", h.Name, net.IP(ip[:]).String())
		queuedPacket.Header.DstMAC = mac
		h.Nics[queuedPacket.Header.SrcMAC].OutBound <- queuedPacket
		delete(h.packetQueue, ip) // Remove from queue
	}
}
