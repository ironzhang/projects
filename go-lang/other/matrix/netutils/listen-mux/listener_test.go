package listen_mux

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func serve(ln net.Listener) {
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go io.Copy(os.Stdout, c)
	}
}

func PrintToServer(addr string) error {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer c.Close()
	fmt.Fprintln(c, addr)
	return nil
}

func TestListen(t *testing.T) {
	addrs := []string{":3000", ":3001", ":3002"}
	ln, err := Listen("tcp", addrs, 0)
	if err != nil {
		t.Fatal(err)
	}
	go serve(ln)

	for _, addr := range addrs {
		if err = PrintToServer(addr); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(100 * time.Millisecond)
	ln.Close()
	time.Sleep(100 * time.Millisecond)
}
