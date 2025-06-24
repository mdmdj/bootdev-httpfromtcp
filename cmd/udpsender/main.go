package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("udp sender startup")

	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	conn, err := net.DialUDP(addr.Network(), nil, addr)

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		inputString, err := inputReader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		n, err := conn.Write([]byte(inputString))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("wrote %d bytes\n", n)
	}

}
