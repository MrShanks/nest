package main

import (
	"time"
)

func main() {
	s := NewSwitch("switch-1")

	A := NewHost("host-A")
	B := NewHost("host-B")

	macA := MAC{0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA}
	ipA := IP{192, 168, 1, 10}
	A.ConfigureNIC(macA, ipA)

	macB := MAC{0xBB, 0xBB, 0xBB, 0xBB, 0xBB, 0xBB}
	ipB := IP{192, 168, 1, 20}
	B.ConfigureNIC(macB, ipB)

	A.Nics[macA].Listen()
	B.Nics[macB].Listen()

	s.Plug(A.Nics[macA])
	s.Plug(B.Nics[macB])

	// Give goroutines time to start up
	time.Sleep(50 * time.Millisecond)

	// B sends to A. B does not know A's MAC address.
	// Expect: B sends ARP, A replies, B sends queued packet, A receives data.
	B.Send([]byte("Hello Host A!"), ipA, macB)

	// Wait for the async operations to complete
	time.Sleep(1 * time.Second)
	B.Send([]byte("Hello again Host A!"), IP{192, 168, 1, 5}, macB)

	time.Sleep(1 * time.Second)
}
