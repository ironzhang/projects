package main

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	os.Exit(0)
}

func TestStringReplace(t *testing.T) {
	s := strings.Replace("dc_etcd_gz01", "etcd", "registry", 1)
	fmt.Printf("%s\n", s)
}
