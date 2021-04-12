package main

import (
	"bufio"
	"log"
	"os"
	"time"
)

func main() {
	f, err := os.Open("./a.pipe")
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Printf("read line: %v", err)
			time.Sleep(time.Second)
			continue
		}
		log.Printf("%s", line)
	}
}
