package main

import (
	"fmt"
	"time"
)

type Frame struct {
	SourceMAC      string
	DestinationMAC string
	Payload        []byte
}

func main() {
	fmt.Println("--- Layer 2 Network Simulator ---")

	// 1. Create network devices
	switch1 := NewSwitch("Switch1")
	switch2 := NewSwitch("Switch2")

	hostA := NewHost("HostA", 1)
	hostB := NewHost("HostB", 1)
	hostC := NewHost("HostC", 2) // HostC has two network cards

	// 2. Connect devices
	// HostA -> Switch1
	nicA := hostA.Interfaces[0]
	portS1A := switch1.ConfigureNewPort()
	nicA.Connect(portS1A)
	go func() {
		for frame := range portS1A {
			nicA.inbound <- frame
		}
	}()

	// HostB -> Switch2
	nicB := hostB.Interfaces[0]
	portS2B := switch2.ConfigureNewPort()
	nicB.Connect(portS2B)
	go func() {
		for frame := range portS2B {
			nicB.inbound <- frame
		}
	}()

	// HostC -> Switch1 (with its first NIC)
	nicC1 := hostC.Interfaces[0]
	portS1C := switch1.ConfigureNewPort()
	nicC1.Connect(portS1C)
	go func() {
		for frame := range portS1C {
			nicC1.inbound <- frame
		}
	}()

	// HostC -> Switch2 (with its second NIC)
	nicC2 := hostC.Interfaces[1]
	portS2C := switch2.ConfigureNewPort()
	nicC2.Connect(portS2C)
	go func() {
		for frame := range portS2C {
			nicC2.inbound <- frame
		}
	}()

	// Connect Switch1 and Switch2 together
	portS1S2 := switch1.ConfigureNewPort()
	portS2S1 := switch2.ConfigureNewPort()
	// This creates a link between the two switches
	go func() {
		for frame := range portS1S2 {
			portS2S1 <- frame
		}
	}()
	go func() {
		for frame := range portS2S1 {
			portS1S2 <- frame
		}
	}()

	// Let the network stabilize and devices start
	time.Sleep(100 * time.Millisecond)

	// 3. Simulate traffic
	fmt.Println("\n--- Simulating Traffic ---")

	// HostA sends a message to HostB.
	// The switches don't know HostB's location yet, so this will be flooded.
	fmt.Println("\n>>> HostA sending to HostB (should flood)")
	hostA.SendData(0, hostB.Interfaces[0].MACaddr, []byte("Hello HostB!"))
	time.Sleep(500 * time.Millisecond)

	// HostB replies to HostA.
	// Switch2 knows HostA is through its connection to Switch1.
	// Switch1 knows HostA's exact port. The frame should be forwarded, not flooded on Switch1.
	fmt.Println("\n>>> HostB replying to HostA (should be more direct)")
	hostB.SendData(0, hostA.Interfaces[0].MACaddr, []byte("Hi HostA, got your message!"))
	time.Sleep(500 * time.Millisecond)

	// HostA sends a message to HostC's first NIC.
	fmt.Println("\n>>> HostA sending to HostC's first NIC")
	hostA.SendData(0, hostC.Interfaces[0].MACaddr, []byte("Hi HostC, from A!"))
	time.Sleep(500 * time.Millisecond)

	// 4. Shutdown
	fmt.Println("\n--- Shutting down simulator ---")
	hostA.Stop()
	hostB.Stop()
	hostC.Stop()
	switch1.Stop()
	switch2.Stop()
	fmt.Println("--- Simulation Finished ---")
}
