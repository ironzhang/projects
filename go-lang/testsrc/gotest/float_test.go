package main

import (
	"fmt"
	"testing"
)

func TestFloat(t *testing.T) {
	f := 1.0
	a := f / 3
	b := a + a + a
	fmt.Printf("a=%f, b=%.9f", a, b)
}
