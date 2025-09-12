package main

import (
	"fmt"
	"time"
)

type Device interface {
	Receive(chan Packet)
	Send([]byte, chan Packet)
}

type Host struct {
}

func NewHost() *Host {
	return &Host{}
}

// Receive reads packets on a channel(wire) and assembles them to reconstruct the message
func (h *Host) Receive(channel chan Packet) {
	var msg = []byte{}
	for p := range channel {
		msg = append(msg, p.payload[:]...)
		fmt.Println("receiving:", string(p.payload))
	}
	fmt.Println("Received message is", string(msg))
}

// Send splits a message into packets chuncks and send them trough the channel(wire)
func (h *Host) Send(message []byte, channel chan Packet) {
	defer close(channel)

	const chunckSize = 25

	for i := 0; i < len(message); i += chunckSize {
		end := min(i+chunckSize, len(message))

		packet := Packet{
			payload: message[i:end],
		}

		channel <- packet

		fmt.Println("Sending:", string(packet.payload))
		time.Sleep(1000 * time.Millisecond)
	}
}
