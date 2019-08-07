package main

import (
	"fmt"
	"net"
	"net/rpc"
)

func NewArithServer() (*rpc.Server, error) {
	svr := rpc.NewServer()
	if err := svr.Register(&Arith{}); err != nil {
		return nil, err
	}
	return svr, nil
}

func main() {
	svr, err := NewArithServer()
	if err != nil {
		fmt.Printf("new arith server: %v\n", err)
		return
	}
	ln, err := net.Listen("tcp", ":2000")
	if err != nil {
		fmt.Printf("net listen: %v\n", err)
		return
	}
	svr.Accept(ln)
}
