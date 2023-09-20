package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: 10000})
	if err != nil {
		log.Printf("dial tcp: %v", err)
		return
	}
	handle(conn)
	time.Sleep(time.Second)
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
	closeR(c)
	defer closeW(c)

	n, err := io.Copy(c, os.Stdin)
	if err != nil {
		log.Printf("copy %d bytes to conn: %v", n, err)
	} else {
		log.Printf("copy %d bytes to conn", n)
	}
}
