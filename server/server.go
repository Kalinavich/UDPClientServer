package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	arg := os.Args
	if len(arg) == 1 {
		fmt.Println("Please provide a port number")
		return
	}

	port := ":" + arg[1]

	address, err := net.ResolveUDPAddr("udp4", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp4", address)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", conn.LocalAddr().String())

	defer conn.Close()

	file, err := os.Open("resources.txt")

	if err != nil {
		fmt.Println("Unable to open resources.txt:", err)
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	resources := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println(err)
			return
		}

		keyValue := strings.SplitN(line, " ", 2)
		if len(keyValue) != 2 {
			continue
		}

		resources[keyValue[0]] = keyValue[1]
	}

	if len(resources) == 0 {
		fmt.Println("Resources empty")
		return
	}

	fmt.Println("Waiting for clients to connect. Server port " + port)

	buf := make([]byte, 64)

	for {
		n, addr, err := conn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error: ", err) //Ошибка подключения
			continue
		}

		go process(conn, addr, string(buf[:n]), &resources)
	}
}

func process(conn *net.UDPConn, addr *net.UDPAddr, key string, resources *map[string]string) {

	fmt.Println("Received ", key, "from ", addr)

	value, find := (*resources)[key]
	fmt.Printf("%q", key)
	fmt.Printf("\n")
	fmt.Printf("%q", value)
	if !find {
		conn.WriteToUDP([]byte("-ERROR-\nResource is not found\n-END-"), addr)
		return
	}

	response := "-BEGIN-\n" + value + "-END-"

	conn.WriteToUDP([]byte(response), addr)
}
