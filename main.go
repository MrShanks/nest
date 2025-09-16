package main

type Packet struct {
	SrcMAC  [6]byte
	DstMAC  [6]byte
	Payload []byte
}

func main() {

	h := NewHost("host-A")
	eth0 := NewNIC("eth0", [6]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, [4]byte{192, 168, 1, 2})

	h.AddNic(eth0)
	go h.Boot()

	s := NewSwitch("switch-1")

	go s.Boot()

	select {}
}
