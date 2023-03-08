package effect

import (
	"fmt"
	"sync"
	"testing"
)

func Test1(t *testing.T) {
	var ints = []int{0, 1, 2, 3}

	var wg sync.WaitGroup
	for _, i := range ints {
		wg.Add(1)
		go func(p *int) {
			defer wg.Done()
			fmt.Println(*p)
		}(&i)
	}
	wg.Wait()
}
