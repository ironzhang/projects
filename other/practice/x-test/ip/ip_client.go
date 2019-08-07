package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
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

func main() {
	var opts Options
	opts.Parse()

	conn, err := net.Dial(opts.Network, opts.Address)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	go mustCopy(os.Stdout, conn)
	mustCopy(conn, os.Stdin)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
