package main

import "fmt"

type Host struct {
	Name string
}

func main() {
	for {
		fmt.Println("[l]ist nics, [a]dd nic, [e]dit nic, [d]elete nic, [q]uit")
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println(err)
		}

		switch input {
		case "l":
			fmt.Println("listing nics")
		case "a":
			fmt.Println("adding nic")
		case "e":
			fmt.Println("editing nic")
		case "d":
			fmt.Println("deleting nic")
		case "1":
			fmt.Println("quitting nic")
		default:
			fmt.Println("please press one of the allowed keys")
			continue
		}
	}
}
