package main

import (
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	conn, err := net.Dial("tcp", "www.baidu.com:443")
	if err != nil {
		t.Fatalf("net dial: %v", err)
	}
	conn.Close()
}
