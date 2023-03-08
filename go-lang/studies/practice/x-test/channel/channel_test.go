package channel

import (
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	done := make(chan struct{})
	afterFunc := func() {
		select {
		case done <- struct{}{}:
		default:
		}
		fmt.Println("after func")
	}

	time.AfterFunc(time.Second, afterFunc)
	time.AfterFunc(2*time.Second, afterFunc)

	<-done
	fmt.Println("done")

	time.Sleep(3 * time.Second)
}

func TestPushChannel(t *testing.T) {
	go func() {
		var done chan struct{}
		done <- struct{}{}
	}()

	time.Sleep(11 * time.Minute)
}
