package ip

import (
	"fmt"
	"log"
	"net"
	"testing"
)

func TestIP(t *testing.T) {
	var err error
	if _, err = net.ResolveIPAddr("ip", "172.16.125.170"); err != nil {
		t.Fatal(err)
	}
	if _, err = net.ResolveIPAddr("ip", "fe80::c27c:d1ff:fea0:4213"); err != nil {
		t.Fatal(err)
	}
	if _, err = net.ResolveIPAddr("ip", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if _, err = net.ResolveIPAddr("ip", "::1"); err != nil {
		t.Fatal(err)
	}
}

func Serve(network, address string) {
	ln, err := net.Listen(network, address)
	if err != nil {
		panic(fmt.Sprintf("server(%s://%s): dial: %v", network, address, err))
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("server: accept: %v", err)
				return
			}
			go HandleConn(network, address, conn)
		}
	}()
}

func HandleConn(network, address string, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("server(%s://%s): read: %v", network, address, err)
			return
		}
		log.Printf("server(%s://%s): %s-->%s: %s", network, address, conn.RemoteAddr().String(), conn.LocalAddr().String(), buf[:n])

		if _, err = conn.Write(buf[:n]); err != nil {
			log.Printf("server(%s://%s): write: %v", network, address, err)
			return
		}
	}
}

func Echo(network, address string) {
	conn, err := net.Dial(network, address)
	if err != nil {
		panic(fmt.Sprintf("client(%s://%s): dial: %v", network, address, err))
	}

	if _, err = conn.Write([]byte("hello")); err != nil {
		log.Printf("client(%s://%s): write: %v", network, address, err)
		return
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("client(%s://%s): read: %v", network, address, err)
		return
	}
	log.Printf("client(%s://%s): %s-->%s: %s", network, address, conn.RemoteAddr().String(), conn.LocalAddr().String(), buf[:n])
}

func TestTPCV4(t *testing.T) {
	Serve("tcp4", ":2000")
	Echo("tcp4", ":2000")
	Echo("tcp", ":2000")
	Echo("tcp", "127.0.0.1:2000")
	Echo("tcp", "localhost:2000")
	Echo("tcp", "172.16.125.41:2000")
	//Echo("tcp", "[::1]:2000")
}

func TestTPCV6(t *testing.T) {
	Serve("tcp6", ":2001")
	Echo("tcp6", ":2001")
	Echo("tcp6", "[::1]:2001")
	//Echo("tcp6", "[fe80::c27c:d1ff:fea0:4213]:2001")
}

func TestTCPALL(t *testing.T) {
	Serve("tcp", ":2002")
	Echo("tcp4", ":2002")
	Echo("tcp6", ":2002")
	Echo("tcp", "127.0.0.1:2002")
	Echo("tcp", "[::1]:2002")
	//Echo("tcp", "[fe80::c27c:d1ff:fea0:4213]:2002")
}
