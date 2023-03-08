package main

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/ironzhang/matrix/netutils/gorpc"
)

func main() {
	c, err := gorpc.Dial("tcp", "localhost:2000")
	if err != nil {
		fmt.Printf("gorpc dial: %v\n", err)
		return
	}
	defer c.Close()

	done := make(chan *rpc.Call, 10)
	go func() {
		for {
			a, b, r := 1, 2, 0
			c.Go("Arith.Add", &Args{A: a, B: b}, &r, done)
			time.Sleep(time.Second)
		}
	}()
	for {
		call := <-done
		if call.Error != nil {
			fmt.Printf("client go: %v\n", call.Error)
		} else {
			args := call.Args.(*Args)
			reply := call.Reply.(*int)
			fmt.Printf("%d + %d = %d\n", args.A, args.B, *reply)
		}
	}
}
