package main

import (
	"log"
	"net"
)

func main() {
	const addr = ":10000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen on %q: %v", addr, err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		conn.Close()
	}
}
