package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"time"

	"github.com/ironzhang/practice/x-test/udp/proto"
)

func ListenUDP(network, address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP(network, addr)
}

func WriteMessage(conn *net.UDPConn, addr net.Addr, msg proto.Message) {
	msg.ServerSend = time.Now()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal: %v", err)
		return
	}
	if _, err = conn.WriteTo(data, addr); err != nil {
		log.Printf("write to: %v", err)
		return
	}
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":2000", "addr")
	flag.Parse()

	ln, err := ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("listen udp: %v", err)
		return
	}
	log.Printf("listen udp on %s", addr)

	ln.SetReadBuffer(4 * 1024 * 1024)
	ln.SetWriteBuffer(4 * 1024 * 1024)

	buf := make([]byte, 1500)
	for {
		n, addr, err := ln.ReadFrom(buf)
		if err != nil {
			log.Printf("read from: %v", err)
			continue
		}
		var msg proto.Message
		if err = json.Unmarshal(buf[:n], &msg); err != nil {
			log.Printf("unmarshal: %v", err)
			continue
		}
		msg.ServerRecv = time.Now()

		go WriteMessage(ln, addr, msg)
	}
}
