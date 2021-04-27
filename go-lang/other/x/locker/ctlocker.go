package flight

import (
	"errors"
	"time"
)

type Locker struct {
	c chan struct{}
}

func (l *Locker) Init(max int) {
	if max <= 0 {
		max = 1
	}
	l.c = make(chan struct{}, max)
}

func (l *Locker) Close() error {
	close(l.c)
	return nil
}

func (l *Locker) Lock(timeout time.Duration) error {
	var c <-chan time.Time
	if timeout > 0 {
		t := time.NewTimer(timeout)
		defer t.Stop()
		c = t.C
	}

	select {
	case l.c <- struct{}{}:
		return nil
	case <-c:
		return errors.New("timeout")
	}
}

func (l *Locker) Unlock(timeout time.Duration) error {
	var c <-chan time.Time
	if timeout > 0 {
		t := time.NewTimer(timeout)
		defer t.Stop()
		c = t.C
	}

	select {
	case <-l.c:
		return nil
	case <-c:
		return errors.New("timeout")
	}
}
