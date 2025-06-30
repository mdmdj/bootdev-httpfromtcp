package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mdmdj/bootdev-httpfromtcp/internal/request"
)

/*
func readLines(f io.ReadCloser, outChan chan string) {
	defer close(outChan)

	readBuffer := make([]byte, 8)

	currentLine := strings.Builder{}

	for {
		bytesRead, err := f.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading file")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for i := 0; i < bytesRead; i++ {
			if readBuffer[i] == '\n' {
				// Complete line found
				outChan <- currentLine.String()
				currentLine.Reset()
			} else {
				// Append to the current line
				currentLine.WriteByte(readBuffer[i])
			}
		}
	}

	if currentLine.String() != "" {
		outChan <- currentLine.String()
	}

} */

func main() {
	//file, err := os.Open("messages.txt")
	fmt.Println("tcp listener startup")

	l, err := net.Listen("tcp", ":42069")

	if err != nil {
		fmt.Println("Error starting server")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connection accepted")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Println("Error reading request")
			log.Println(err.Error())
			continue
		}

		rl := request.RequestLine
		if rl.Method == "" {
			log.Println("Error reading request line:", rl)
			continue
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", rl.Method)
		fmt.Printf("- Target: %s\n", rl.RequestTarget)
		fmt.Printf("- Version: %s\n", rl.HttpVersion)

		fmt.Println("Connection closed")
	}

	//os.Exit(0)

}
