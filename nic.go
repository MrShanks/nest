package main

import (
	"fmt"
	"net"
)

type NIC struct {
	Mac      MAC
	Ip       IP
	OutBound chan *Packet
	InBound  chan *Packet
	Host     *Host
}

func (n *NIC) Listen() {
	go func() {
		for packet := range n.InBound {
			isBroadcast := packet.Header.DstMAC == MAC{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
			isForMe := packet.Payload.DstIP == n.Ip

			// --- ARP Request Handling ---
			if isBroadcast && isForMe {
				fmt.Printf("[%s] NIC [%s] received an ARP request for me\n", n.Host.Name, net.HardwareAddr(n.Mac[:]).String())

				// Learn the sender's info from their request
				n.Host.updateArpCache(packet.Payload.SrcIP, packet.Header.SrcMAC)

				// Craft the ARP reply
				replyPacket := &Packet{
					Header: EthernetHeader{
						SrcMAC: n.Mac,
						DstMAC: packet.Header.SrcMAC, // Send it directly back to the requester
					},
					Payload: IPv4Packet{
						SrcIP: n.Ip,
						DstIP: packet.Payload.SrcIP,
						Data:  []byte("ARP_REPLY"),
					},
				}
				// Send the reply
				n.OutBound <- replyPacket
				continue // Done with this packet
			}

			// --- ARP Reply Handling ---
			if !isBroadcast && isForMe && string(packet.Payload.Data) == "ARP_REPLY" {
				fmt.Printf("[%s] NIC [%s] received an ARP reply\n", n.Host.Name, net.HardwareAddr(n.Mac[:]).String())
				// Update our cache with the information from the reply
				n.Host.updateArpCache(packet.Payload.SrcIP, packet.Header.SrcMAC)
				continue
			}

			// --- Regular Data Handling ---
			if isForMe {
				fmt.Printf("[%s] NIC [%s] received data: '%s'\n", n.Host.Name, net.HardwareAddr(n.Mac[:]).String(), string(packet.Payload.Data))
			}
		}
	}()
}
