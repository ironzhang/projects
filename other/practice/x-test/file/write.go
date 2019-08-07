package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Create("./a.pipe")
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("stat: %v", err)
	}
	fmt.Fprintf(f, "modify at %s\n", fi.ModTime())
}
