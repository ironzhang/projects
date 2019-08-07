package print

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPrint1(t *testing.T) {
	for i := 0; i <= 20; i++ {
		fmt.Printf("%010d\n", i)
	}
}

func TestPrint2(t *testing.T) {
	fmt.Print(1, 2, 3, 1, "2", 3.0, "\n")
	fmt.Print(fmt.Sprint(1, 2, 3, 1, "2", 3.0, "\n"))
}

func TestPrintBin(t *testing.T) {
	s1 := "hello"
	s2 := string([]byte{1, 2, 3})
	s3 := s1 + s2

	fmt.Printf("s1: %s\n", s1)
	fmt.Printf("s2: %s\n", s2)
	fmt.Printf("s3: %s\n", s3)

	v := struct {
		A string
		B string
		C string
	}{
		A: s1,
		B: s2,
		C: s3,
	}

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("json: %s\n", data)
}
