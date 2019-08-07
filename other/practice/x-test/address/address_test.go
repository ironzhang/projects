package main

import (
	"log"
	"net"
	"testing"
	"time"
)

func RunServer(network, address string) {
	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("Addr: %s\n", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			return
		}
		log.Printf("Server: LocalAddr: %s, RemoteAddr: %s", conn.LocalAddr(), conn.RemoteAddr())
		conn.Close()
	}
}

func ConnectServer(network, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	log.Printf("Client: LocalAddr: %s, RemoteAddr: %s", conn.LocalAddr(), conn.RemoteAddr())
	conn.Close()
	return nil
}

func TestAddress(t *testing.T) {
	go RunServer("tcp4", ":7000")
	time.Sleep(time.Second)
	if err := ConnectServer("tcp", ":7000"); err != nil {
		t.Fatalf("connect server: %v", err)
	}
	if err := ConnectServer("tcp", "127.0.0.1:7000"); err != nil {
		t.Fatalf("connect server: %v", err)
	}
	if err := ConnectServer("tcp", "172.16.127.243:7000"); err != nil {
		t.Fatalf("connect server: %v", err)
	}
	time.Sleep(time.Second)
}
