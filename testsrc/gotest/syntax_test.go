package main

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"unsafe"
)

func TestSliceRange(t *testing.T) {
	ss := []string{"a", "b", "c", "d"}
	for i, s := range ss {
		fmt.Println(s)
		if i == 0 {
			fmt.Println("append e")
			ss = append(ss, "e")
		}
	}
}

func TestSliceForeach(t *testing.T) {
	ss := []string{"a", "b", "c", "d"}
	for i := 0; i < len(ss); i++ {
		fmt.Println(ss[i])
		if i == 0 {
			fmt.Println("append e")
			ss = append(ss, "e")
		}
	}
}

func TestIntMod(t *testing.T) {
	fmt.Printf("%d\n", 5%3)
	fmt.Printf("%d\n", -5%3)
}

func TestIntMat(t *testing.T) {
	fmt.Printf("%d\n", math.MaxInt16)
}

func TestSizeof(t *testing.T) {
	var a interface{}
	var ch chan interface{}
	fmt.Printf("interface{}: %d\n", unsafe.Sizeof(a))
	fmt.Printf("chan: %d\n", unsafe.Sizeof(ch))
}

func TestPanic(t *testing.T) {
	//panic("panic")
}

func TestDiv(t *testing.T) {
	i := 2
	fmt.Printf("1.0/2=%f\n", float64(1.0/float64(i)))
}

func testDeferError() (err error) {
	defer func() {
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	}()
	return errors.New("return defer error")
}

func TestDeferError(t *testing.T) {
	testDeferError()
}
