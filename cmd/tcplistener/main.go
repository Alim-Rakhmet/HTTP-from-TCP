package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		line := ""

		buf := make([]byte, 8)
		for {
			n, err := f.Read(buf)

			if err != nil {
				if errors.Is(err, io.EOF) {
					out <- line
					break
				}
				log.Fatal(err.Error())
			}

			buf = buf[:n]

			newLine := buf
			for {
				sep := bytes.IndexByte(newLine, '\n')
				if sep != -1 {
					line += string(newLine[:sep])
					out <- line
					line = ""
					newLine = newLine[sep+1:]
				} else {
					line += string(newLine)
					break
				}
			}
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	lines := getLinesChannel(conn)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
