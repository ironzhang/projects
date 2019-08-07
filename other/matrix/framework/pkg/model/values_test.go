package model

import (
	"reflect"
	"testing"
)

func TestValues(t *testing.T) {
	var ok bool
	var i int
	var v Values

	_, ok = v.GetValue("name")
	if ok {
		t.Errorf("get value item")
	}

	_, ok = v.GetInterface("name")
	if ok {
		t.Errorf("get interface item")
	}

	if err := v.Register("name", &i); err != nil {
		t.Fatalf("register: %v", err)
	}

	x, ok := v.GetValue("name")
	if !ok {
		t.Errorf("can not get value item")
	}
	if got, want := x.Load(), &i; !reflect.DeepEqual(got, want) {
		t.Errorf("value: got(%v) != want(%v)", got, want)
	}

	y, ok := v.GetInterface("name")
	if !ok {
		t.Errorf("can not get interface item")
	}
	if got, want := y, &i; !reflect.DeepEqual(got, want) {
		t.Errorf("value: got(%v) != want(%v)", got, want)
	}
}
