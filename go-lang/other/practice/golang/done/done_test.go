package done_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/ironzhang/golang/done"
)

func run(ctx *done.Context, key interface{}) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%v goroutine done\n", key)
			ctx.OK()
			return
		}
	}
}

func TestDone(t *testing.T) {
	//runtime.GOMAXPROCS(1)
	var dg done.DoneGroup
	for i := 0; i < 10; i++ {
		c, err := dg.Add(i)
		if err != nil {
			t.Fatalf("go failed, i[%d]", i)
		}
		go run(c, i)
	}
	for i := 0; i < 10; i++ {
		if err := dg.Done(i, true); err != nil {
			t.Errorf("done failed, i[%d]", i)
		}
	}
}

func TestDoneAll(t *testing.T) {
	runtime.GOMAXPROCS(1)
	var dg done.DoneGroup
	for i := 0; i < 10; i++ {
		c, err := dg.Add(i)
		if err != nil {
			t.Fatalf("go failed, i[%d]", i)
		}
		go run(c, i)
	}
	dg.DoneAll(true)
}
