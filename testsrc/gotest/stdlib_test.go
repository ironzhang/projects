package main

import (
	"context"
	"fmt"
	"testing"
)

func TestRelovse(t *testing.T) {
	//net.LookupNS()
}

func TestPrintContext(t *testing.T) {
	ctx0 := context.Background()
	ctx1 := context.WithValue(ctx0, "key", "value")
	fmt.Printf("%v\n", ctx0)
	fmt.Printf("%v\n", ctx1)
}
