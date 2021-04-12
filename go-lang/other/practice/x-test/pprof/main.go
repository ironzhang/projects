package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

func runTask() {
	for {
		sum := 0
		for i := 0; i < 1000000000; i++ {
			sum += i
		}
		fmt.Println(sum)
		time.Sleep(time.Second)
	}
}

func run(n int) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runTask()
		}()
	}
	wg.Wait()
}

func main() {
	go run(10)

	var _ = pprof.Index
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
