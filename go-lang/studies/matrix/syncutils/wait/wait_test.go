package wait

import (
	"fmt"
	"testing"
	"time"
)

func TestWaiter(t *testing.T) {
	w := waiter{done: make(chan struct{})}
	time.AfterFunc(time.Second, func() { w.Done(nil, ErrTimeout) })
	_, err := w.Wait()
	if got, want := err, ErrTimeout; got != want {
		t.Errorf("wait: %v != %v", got, want)
	}
}

func TestGroup(t *testing.T) {
	var g Group
	var ts []Token
	for i := 0; i < 10; i++ {
		k := fmt.Sprint(i)
		tk, err := g.Add(k, time.Duration(i)*time.Second)
		if err != nil {
			t.Fatalf("add: %v", err)
		}
		ts = append(ts, tk)
		go g.Done(k, i, nil)
	}
	for i, tk := range ts {
		v, err := tk.Wait()
		if err != nil {
			t.Errorf("wait: %v", err)
			continue
		}
		got, want := v.(int), i
		if got != want {
			t.Errorf("value: %v != %v", got, want)
			continue
		}
		t.Logf("%dth token wait success: value=%d", i, got)
	}
}

func TestGroupTimeout(t *testing.T) {
	var g Group
	tk, err := g.Add("1", time.Second)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	_, err = tk.Wait()
	if got, want := err, ErrTimeout; got != want {
		t.Errorf("wait: %v != %v", got, want)
	}
}
