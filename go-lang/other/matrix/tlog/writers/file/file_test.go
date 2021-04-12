package file

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFile0(t *testing.T) {
	f, err := Open(Options{
		Dir:      os.TempDir(),
		Program:  "TestFile0",
		MaxBytes: 500,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < 10; i++ {
		fmt.Fprintf(f, "%d: Hello, world\n", i)
	}
}

func TestFile1(t *testing.T) {
	f, err := Open(Options{
		Dir:     os.TempDir(),
		Program: "TestFile1",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprintln(f, time.Now())
	time.Sleep(flushDelay + time.Second)
	fmt.Fprintln(f, time.Now())
	fmt.Fprintln(f, time.Now())
	time.Sleep(flushDelay + time.Second)
}
