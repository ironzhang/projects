package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 10000})
	if err != nil {
		log.Fatalf("listen tcp: %v", err)
		return
	}
	log.Printf("listen tcp on %q", ln.Addr().String())

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handle(conn)
	}
}

func closeR(c *net.TCPConn) {
	log.Printf("close conn read pipe")
	c.CloseRead()
}

func closeW(c *net.TCPConn) {
	log.Printf("close conn write pipe")
	c.CloseWrite()
}

func handle(c *net.TCPConn) {
	closeW(c)
	defer closeR(c)

	if n, err := io.Copy(os.Stdout, c); err != nil {
		log.Printf("copy %d bytes from conn: %v", n, err)
	} else {
		log.Printf("copy %d bytes from conn", n)
	}
}
