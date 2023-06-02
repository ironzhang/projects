package filewatch

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	path := "./testdata/test_watcher/test.txt"

	w := NewWatcher(100 * time.Millisecond)
	defer w.Stop()

	n := 0
	w.WatchFile(context.Background(), path, func(path string) bool {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}

		fmt.Printf("n=%d, data=%q\n", n, data)
		n++

		return false
	})

	for i := 0; i < 4; i++ {
		WriteTestFile(path, fmt.Sprint(i))
		time.Sleep(150 * time.Millisecond)
	}
}
