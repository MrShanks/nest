package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Host struct {
	Name string
	Nics []*NIC
	wg   *sync.WaitGroup
}

func NewHost(name string) *Host {
	return &Host{
		Name: name,
		Nics: make([]*NIC, 0),
		wg:   &sync.WaitGroup{},
	}
}

func (h *Host) AddNic(nic *NIC) {
	h.Nics = append(h.Nics, nic)
}

func (h *Host) Boot() {
	for _, nic := range h.Nics {
		h.StartListenerOn(nic)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		// Wait for the shutdown signal
		<-shutdown

		fmt.Println("\nShutting down gracefully...")

		for _, nic := range h.Nics {
			nic.Socket.Close()
		}
	}()

	h.wg.Wait()

	for _, nic := range h.Nics {
		os.Remove(nic.SocketPath)
	}
}

func (h *Host) StartListenerOn(nic *NIC) {
	nic.SocketPath = fmt.Sprintf("./%s_%s.sock", h.Name, nic.Name)

	// Clean up old dangling sockets
	os.Remove(nic.SocketPath)

	fmt.Printf("[%s] listening on [%s]\n", h.Name, nic.Name)

	var err error
	nic.Socket, err = net.Listen("unix", nic.SocketPath)
	if err != nil {
		log.Fatal(err)
	}

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		for {
			conn, err := nic.Socket.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					fmt.Printf("[%s] listener on [%s] is now closed.\n", h.Name, nic.Name)
					return
				}
				log.Printf("%v\n", err)
				return
			}

			go handleConnection(conn, nic)
		}
	}()
}

func handleConnection(conn net.Conn, nic *NIC) {
	defer conn.Close()

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Received message on [%s]: %s", nic.Name, buf[:n])

}

func (h *Host) GetNicByName(nicName string) *NIC {
	var nic *NIC
	for _, n := range h.Nics {
		if nic.Name == nicName {
			nic = n
		}
	}
	return nic
}

func (h *Host) ConnectTo(nicName string, s *Switch) {
	nic := h.GetNicByName(nicName)
	s.Ports = append(s.Ports, nic.Socket)
}
