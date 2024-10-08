package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)


const timeout = 5 * time.Second

type server struct {
	address string
	active  bool
}

var servers []*server = []*server{
	{address: "127.0.0.1:8081", active: true},
	{address: "127.0.0.1:8082", active: true},
	{address: "127.0.0.1:8083", active: true},
}

func (s *server) activate() {
	s.active = true
}

func (s *server) deactivate() {
	s.active = false
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stdout, "Listening for connections on 127.0.0.1:8080...")

	go checkHealthyServers()
	acceptRequests(ln)
}

func acceptRequests(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
		fmt.Println()
	}
}

var serverPos int = -1

func getNextServer() (*server, error) {
	inactiveCount := 0

	serverPos = (serverPos + 1) % len(servers)
	for !servers[serverPos].active {
		serverPos = (serverPos + 1) % len(servers)
		inactiveCount++

		if inactiveCount == len(servers) {
			return nil, errors.New("all servers are down")
		}
	}

	return servers[serverPos], nil
}

func handleConnection(conn net.Conn) {
	fmt.Fprintf(os.Stdout, "Received request from %s\n", conn.RemoteAddr())
	defer conn.Close()

	clientRes, err := readFromConnection(conn)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	retries := 3
	for retries > 0 {
		srv, err := getNextServer()
		if err != nil {
			buf := bytes.Buffer{}
			buf.WriteString("HTTP/1.1 502 Bad Gateway\r\n")
			buf.WriteString("\r\n")
			conn.Write(buf.Bytes())
			conn.Close()
			return
		}

		beConn, err := net.DialTimeout("tcp", srv.address, timeout)
		if err != nil {
			log.Println(err)
			srv.deactivate()
			retries--
			continue
		}

		_, err = beConn.Write([]byte(clientRes))
		if err != nil {
			log.Println(err)
			srv.deactivate()
			beConn.Close()
			retries--
			continue
		}

		s := fmt.Sprintf("Response from server %s: ", srv.address)
		backendRes, err := readFromConnection(beConn)
		if err != nil {
			log.Println(err)
			srv.deactivate()
			beConn.Close()
			retries--
			continue
		}
		fmt.Fprint(os.Stdout, s+backendRes)
		conn.Write([]byte(backendRes))
		break
	}
}

func readFromConnection(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	reader := bufio.NewReader(conn)

	buf := bytes.Buffer{}
	contentLength := 0
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		buf.WriteString(s)

		if strings.HasPrefix(s, "Content-Length:") {
			lengthInStr := strings.Split(s, ":")[1]
			contentLength, err = strconv.Atoi(strings.TrimSpace(lengthInStr))
			if err != nil {
				log.Fatal(err)
			}
		}

		if s == "\r\n" {
			break
		}
	}

	for contentLength != 0 {
		b, err := reader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		buf.WriteByte(b)
		contentLength--
	}

	return buf.String(), nil
}

func checkHealthyServers() {
	for {
		time.Sleep(10 * time.Second)
		for _, server := range servers {
			if isHealthy(server.address) {
				server.activate()
			} else {
				server.deactivate()
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func isHealthy(serverAddress string) bool {
	beConn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Printf("Could not connected to server %s", serverAddress)
		return false
	}

	buf := bytes.Buffer{}
	buf.WriteString("GET /health HTTP/1.1\r\n")
	buf.WriteString("\r\n")
	beConn.Write(buf.Bytes())

	s := fmt.Sprintf("Response from Health Check in server %s: ", serverAddress)
	res, err := readFromConnection(beConn)
	if err != nil {
		return false
	}

	fmt.Fprint(os.Stdout, s+res)
	beConn.Close()

	tokens := strings.Split(res, " ")
	if tokens[1] != "200" && tokens[1] != "204" {
		return false
	}

	return true
}