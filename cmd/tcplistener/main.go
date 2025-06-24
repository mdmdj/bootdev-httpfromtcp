package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	outChan := make(chan string)

	go readLines(f, outChan)

	return outChan
}

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

}

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

		linesChannel := getLinesChannel(conn)

		for line := range linesChannel {
			fmt.Println(line)
		}

		fmt.Println("Connection closed")
	}

	//os.Exit(0)

}
