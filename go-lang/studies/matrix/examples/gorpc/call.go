package main

import (
	"fmt"
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

	a, b, r := 1, 2, 0
	for {
		if err = c.Call("Arith.Add", &Args{A: a, B: b}, &r); err != nil {
			fmt.Printf("client call: %v\n", err)
		} else {
			fmt.Printf("%d + %d = %d\n", a, b, r)
		}
		time.Sleep(time.Second)
	}
}
