package main

import (
	"fmt"
	"sync"
)

//func init() {
//	fmt.Println("init")
//}

func run(ch chan int) {
	select {
	case ch <- 1:
	default:
	}
}

func main() {
	fmt.Println("main start")

	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		run(ch)
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("main end")
}
