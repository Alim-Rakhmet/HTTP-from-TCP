package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:42069")
	if err != nil {
		log.Fatalf("addresing err: %s", err.Error())
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("dialing err: %s", err.Error())
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Println(err.Error())
		}
	}
}
