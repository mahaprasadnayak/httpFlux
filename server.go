package main

import (
	"flag"
	"fmt"
	"httpFlux/node"
	"log"
	"net"
	"os"
)

func main() {
	var hostname, port string
	flag.StringVar(&hostname, "h", "127.0.0.1", "hostname")
	flag.StringVar(&port, "p", "8081", "port")
	flag.Parse()

	ln, err := net.Listen("tcp", hostname+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "Listening for connections on %s:%s...\n", hostname, port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go node.HandleConnection(conn)
		fmt.Println()
	}

}

