package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func randomMAC() (net.HardwareAddr, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	buf[0] = (buf[0] | 2) & 0xfe

	return net.HardwareAddr(buf), nil
}

type Nic struct {
	Name string
	MAC  net.HardwareAddr
	IP   net.IP
}

type Host struct {
	Name string
	Nics []*Nic
}

func NewHost(name string) *Host {
	return &Host{
		Name: name,
	}
}

func NewNic(name string) *Nic {
	mac, err := randomMAC()
	if err != nil {
		fmt.Println(err)
	}

	return &Nic{
		Name: name,
		MAC:  mac,
	}
}

func (h *Host) AddNic(name string) {
	h.Nics = append(h.Nics, NewNic(name))
}

func (h *Host) List() {
	for _, nic := range h.Nics {
		fmt.Println(nic.Name)
	}
}

func (h *Host) PrintNicConf(index int) {
	fmt.Println("Name: ", h.Nics[index].Name)
	fmt.Println("IP  : ", h.Nics[index].IP)
	fmt.Println("MAC : ", h.Nics[index].MAC)
}

type Menu struct {
	options map[string]func()
}

func PopulateMenu() *Menu {
	m := &Menu{
		options: make(map[string]func()),
	}

	m.options["s"] = func() {
		if len(h.Nics) < 1 {
			fmt.Println("There are no nic configured yet")
			return
		}
		fmt.Println("Press number associated with nic to see its configuration")
		for i, nic := range h.Nics {
			fmt.Printf("%2d: %s\n", i, nic.Name)
		}

		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println(err)
		}
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(err)
		}
		h.PrintNicConf(index)
	}
	m.options["l"] = func() {
		fmt.Printf("Listing nics: \n")
		h.List()
	}

	m.options["a"] = func() {
		fmt.Printf("Adding nic:\n")
		fmt.Println("Insert name for Nic: ")
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println(err)
		}

		h.AddNic(input)
	}
	m.options["e"] = func() { fmt.Println("Editing nic") }
	m.options["d"] = func() { fmt.Println("Deleting nic") }
	m.options["q"] = func() {
		fmt.Println("Quitting...")
		os.Exit(0)
	}

	return m
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

var h *Host

func main() {
	clearScreen()
	h = NewHost("Host-A")

	m := &Menu{
		options: map[string]func(){
			"s": func() {
				fmt.Printf("[%s]:~$\n", h.Name)

				cmd := exec.Command("ls")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				cmd.SysProcAttr = &syscall.SysProcAttr{
					Chroot: "./",
				}

				if err := cmd.Run(); err != nil {
					fmt.Println(err)
				}
			},
		},
	}

	m.options["s"]()

	///menu := PopulateMenu()
	///for {
	///	fmt.Printf("%s:~$: [s]how configuration, [l]ist nics, [a]dd nic, [e]dit nic, [d]elete nic, [q]uit\n", h.Name)
	///	var input string
	///	_, err := fmt.Scan(&input)
	///	if err != nil {
	///		fmt.Println(err)
	///	}

	///	if action, ok := menu.options[input]; ok {
	///		action()
	///	}
	///}
}
