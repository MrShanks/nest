package main

import "time"

func main() {

	a := NewHost()
	b := NewHost()

	l := NewLink(a, b)

	l.transmit([]byte("this is the message to send from host a to host b let's see how many it gets splitted let'see if we can make it longer so that it need to be send in multiple chunks it should be a lot lot longer because every chunk is 50 bytes"))

	time.Sleep(20 * time.Second)
}
