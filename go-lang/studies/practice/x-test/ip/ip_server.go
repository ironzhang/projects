package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

type Options struct {
	Network string
	Address string
}

func (o *Options) Parse() {
	flag.StringVar(&o.Network, "net", "tcp", "network")
	flag.StringVar(&o.Address, "addr", ":10000", "address")
	flag.Parse()
}

func ListenAndServe(network, address string) error {
	ln, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go ServeConn(conn)
	}
}

func ServeConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		fmt.Fprintln(c, input.Text())
	}
	c.Close()
}

func main() {
	var opts Options
	opts.Parse()

	if err := ListenAndServe(opts.Network, opts.Address); err != nil {
		log.Printf("listen and serve: %v", err)
	}
}
