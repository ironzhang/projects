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
	//	for i := 0; i < len(ss); i++ {
	//		fmt.Println(ss[i])
	//		if i == 0 {
	//			//fmt.Println("append e")
	//			ss = append(ss, "e")
	//		}
	//	}

	for i, c := range ss {
		fmt.Println(c)
		if i == 0 {
			//fmt.Println("append e")
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
	fmt.Printf("1.0/2=%v\n", float64(1.0/i))
	fmt.Printf("1.0/2.0=%v\n", float64(1.0/float64(i)))
}

func TestMod(t *testing.T) {
	x := -1
	fmt.Printf("-1%%10=%d\n", x%10)
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

func TestChan(t *testing.T) {
	ch := make(chan int)
	ch <- 1
}

func TestNilMap(t *testing.T) {
	var m map[string]string

	m = nil
	fmt.Printf("m[a]=%s\n", m["a"])
	fmt.Printf("m[b]=%s\n", m["b"])
}
