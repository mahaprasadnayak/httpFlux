package node

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func HandleConnection(conn net.Conn) {
	fmt.Fprintf(os.Stdout, "Received request from %s\n", conn.RemoteAddr())
	reader := bufio.NewReader(conn)

	var path string

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if s == "\r\n" {
			break
		}

		tokens := strings.Split(s, " ")
		if tokens[0] == "GET" {
			path = tokens[1]
		}

		fmt.Fprint(os.Stdout, s)
	}

	buf := handleRoute(path)
	conn.Write(buf.Bytes())
	conn.Close()
}

func handleRoute(path string) *bytes.Buffer {
	buf := bytes.Buffer{}

	switch path {
	case "/health":
		buf.Write([]byte("HTTP/1.1 204 No Content\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("\r\n"))
	case "/":
		buf.Write([]byte("HTTP/1.1 200 OK\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("Content-Length: 27\r\n"))
		buf.Write([]byte("\r\n"))
		buf.Write([]byte("Hello From Flux Server !!!!! \r\n"))
	default:
		buf.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("\r\n"))
	}

	return &buf
}