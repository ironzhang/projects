package wait

import (
	"errors"
	"sync"
	"time"
)

var ErrTimeout = errors.New("wait timeout")
var ErrKeyDuplicate = errors.New("key duplicate")

type waiter struct {
	done  chan struct{}
	err   error
	value interface{}
}

func (w *waiter) Wait() (interface{}, error) {
	<-w.done
	return w.value, w.err
}

func (w *waiter) Done(value interface{}, err error) {
	w.err = err
	w.value = value
	close(w.done)
}

type Token interface {
	Wait() (interface{}, error)
}

type Group struct {
	mu      sync.Mutex
	waiters map[interface{}]*waiter
}

func (g *Group) Add(key interface{}, timeout time.Duration) (Token, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	w, ok := g.waiters[key]
	if ok {
		return nil, ErrKeyDuplicate
	}

	w = &waiter{done: make(chan struct{})}
	if g.waiters == nil {
		g.waiters = make(map[interface{}]*waiter)
	}
	g.waiters[key] = w
	if timeout > 0 {
		time.AfterFunc(timeout, func() { g.Done(key, nil, ErrTimeout) })
	}
	return w, nil
}

func (g *Group) Done(key interface{}, value interface{}, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	w, ok := g.waiters[key]
	if !ok {
		return
	}

	delete(g.waiters, key)
	w.Done(value, err)
}
