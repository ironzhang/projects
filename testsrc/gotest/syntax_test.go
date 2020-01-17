package main

import (
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
