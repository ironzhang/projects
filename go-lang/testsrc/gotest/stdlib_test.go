package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
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

func TestTimeHour(t *testing.T) {
	n := time.Now()
	h := n.Hour()
	fmt.Printf("h=%d\n", h)

	t1 := n.Add(-1 * time.Duration(h) * time.Hour)
	fmt.Printf("t1=%s\n", t1)

	t2 := t1.Truncate(time.Hour)
	fmt.Printf("t2=%s\n", t2)
}
