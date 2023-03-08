package collect

import (
	"os"
	"sync"
	"testing"
	"time"
)

func Function1() {
	collect := Default.Collect("Function1")
	defer collect()

	time.Sleep(10 * time.Millisecond)
}

func Function2() {
	collect := Default.Collect("Function2")
	defer collect()

	time.Sleep(20 * time.Millisecond)
}

func TestCollect(t *testing.T) {
	for n := 0; n < 2; n++ {
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Function1()
			}()
		}
		for i := 0; i < 200; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Function2()
			}()
		}
		wg.Wait()
		Default.Print(os.Stdout)
		Default.Clear()
	}
}
