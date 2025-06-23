package main

import (
	"fmt"
	"io"
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
	file, err := os.Open("messages.txt")

	if err != nil {
		fmt.Println("Error opening file")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer file.Close()

	linesChannel := getLinesChannel(file)

	for line := range linesChannel {
		fmt.Printf("read: %s\n", line)
	}

	os.Exit(0)

}
