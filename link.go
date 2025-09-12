package main

type Link struct {
	A, B    Device
	channel chan Packet
}

func NewLink(a, b Device) *Link {
	return &Link{
		A:       a,
		B:       b,
		channel: make(chan Packet),
	}
}

func (l *Link) transmit(msg []byte) {
	go l.A.Send(msg, l.channel)
	go l.B.Receive(l.channel)
}
