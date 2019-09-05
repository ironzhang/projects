package main

import (
	"fmt"
	"testing"
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
