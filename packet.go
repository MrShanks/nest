package main

type MAC [6]byte
type IP [4]byte

type EthernetHeader struct {
	SrcMAC MAC
	DstMAC MAC
}

type IPv4Packet struct {
	SrcIP IP
	DstIP IP
	Data  []byte
}

type Packet struct {
	Header  EthernetHeader
	Payload IPv4Packet
}
