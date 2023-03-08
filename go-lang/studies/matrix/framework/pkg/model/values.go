package model

import (
	"fmt"
	"reflect"
	"sync"
)

type Reloader interface {
	Reload() error
}

type Value struct {
	mu  sync.Mutex
	ptr interface{}
}

func (v *Value) Load() interface{} {
	return v.ptr
}

func (v *Value) Store(a interface{}) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	return setValue(v.ptr, a)
}

func (v *Value) Reload() error {
	if r, ok := v.ptr.(Reloader); ok {
		return r.Reload()
	}
	return nil
}

type Values struct {
	m map[string]*Value
}

func (p *Values) Register(name string, ptr interface{}) error {
	if p.m == nil {
		p.m = make(map[string]*Value)
	}
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("value invalid: %s, %T", name, ptr)
	}
	if _, ok := p.m[name]; ok {
		return fmt.Errorf("value duplicate: %s", name)
	}
	p.m[name] = &Value{ptr: ptr}
	return nil
}

func (p *Values) GetValue(name string) (*Value, bool) {
	if v, ok := p.m[name]; ok {
		return v, true
	}
	return nil, false
}

func (p *Values) GetInterface(name string) (interface{}, bool) {
	if v, ok := p.m[name]; ok {
		return v.ptr, true
	}
	return nil, false
}

func (p *Values) Interfaces() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range p.m {
		m[k] = v.ptr
	}
	return m
}
