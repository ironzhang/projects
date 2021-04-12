package main

import (
	"fmt"
	"syscall"
	"testing"
	"time"
)

func SetSysTime(t time.Time) error {
	tv := syscall.NsecToTimeval(t.UnixNano())
	return syscall.Settimeofday(&tv)
}

func TestTimeBefore(t *testing.T) {
	n1 := time.Now()
	defer func() {
		SetSysTime(n1)
	}()

	t1 := n1.Add(2 * time.Hour)
	t2 := n1.Add(3 * time.Hour)
	if err := SetSysTime(t2); err != nil {
		t.Fatalf("set sys time: %v", err)
	}
	n2 := time.Now()

	fmt.Printf("n1=[%v], n2=[%v]\n", n1, n2)
	fmt.Printf("t1=[%v], t2=[%v]\n", t1, n2)

	if t2.Before(t1) {
		t.Errorf("t2 before t1, t2=[%v][%d], t1=[%v][%d]", t2, t2.UnixNano(), t1, t1.UnixNano())
	}
	if n2.Before(t1) {
		t.Errorf("n2 before t1, n2=[%v][%d], t1=[%v][%d]", n2, n2.UnixNano(), t1, t1.UnixNano())
	}
}
