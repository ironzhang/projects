package done

import (
	"errors"
	"sync"
)

type Context struct {
	done chan struct{}
	ok   chan struct{}
}

func newContext() *Context {
	return &Context{
		done: make(chan struct{}),
		ok:   make(chan struct{}),
	}
}

func (c *Context) Done() <-chan struct{} {
	return c.done
}

func (c *Context) OK() {
	close(c.ok)
}

type DoneGroup struct {
	mu sync.RWMutex
	m  map[interface{}]*Context
}

func (g *DoneGroup) addContext(key interface{}) (*Context, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.m == nil {
		g.m = make(map[interface{}]*Context)
	}
	c, ok := g.m[key]
	if ok {
		return nil, false
	}
	c = newContext()
	g.m[key] = c
	return c, true
}

func (g *DoneGroup) delContext(key interface{}) (*Context, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.m == nil {
		g.m = make(map[interface{}]*Context)
	}
	c, ok := g.m[key]
	if !ok {
		return nil, false
	}
	delete(g.m, key)
	return c, true
}

func (g *DoneGroup) Add(key interface{}) (*Context, error) {
	c, ok := g.addContext(key)
	if !ok {
		return nil, errors.New("key existed")
	}
	return c, nil
}

func (g *DoneGroup) Done(key interface{}, wait bool) error {
	c, ok := g.delContext(key)
	if !ok {
		return errors.New("key not exist")
	}
	close(c.done)
	if wait {
		<-c.ok
	}
	return nil
}

func (g *DoneGroup) DoneAll(wait bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.m == nil {
		g.m = make(map[interface{}]*Context)
	}
	for k, c := range g.m {
		delete(g.m, k)
		close(c.done)
		if wait {
			<-c.ok
		}
	}
}
