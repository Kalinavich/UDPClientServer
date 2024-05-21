package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	arg := os.Args
	if len(arg) != 3 {
		fmt.Println("Please enter the arguments: host name, port number and name request resource, in the form HOST:PORT RESOURCE")
		return
	}

	serverAddr := arg[1]
	resourceMsg := arg[2]

	if len(resourceMsg) > 64 {
		fmt.Fprintln(os.Stderr, "Request number of characters exceeded")
		return
	}

	addr, err := net.ResolveUDPAddr("udp4", serverAddr)

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp4", nil, addr)

	defer conn.Close()

	conn.SetWriteBuffer(64)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write([]byte(resourceMsg))

	if err != nil {
		fmt.Println(resourceMsg, err)
	}

	buf := make([]byte, 1024)

	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error on receiving: ", err)
		return
	}

	message := string(buf[0:n])

	if message[:8] != "-ERROR-\n" {
		fmt.Println(message)
	} else {
		fmt.Fprintln(os.Stderr, message)
	}
}
