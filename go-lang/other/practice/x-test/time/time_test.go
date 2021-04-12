package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()

	fmt.Printf("Unix: %v\n", now.Unix())
	fmt.Printf("UTC.Unix: %v\n", now.UTC().Unix())
}
