package discovery

import (
	"reflect"
	"testing"
)

func TestServiceRefresh0(t *testing.T) {
	s := newService("name")
	kvs := map[string][]byte{
		"a1": []byte("1"),
		"a3": []byte("1"),
		"a2": []byte("1"),
	}
	addrs := []string{"a1", "a2", "a3"}
	s.Refresh(kvs)

	if got, want := s.name, "name"; got != want {
		t.Errorf("name: %q != %q", got, want)
	}
	if got, want := s.addrs, addrs; !reflect.DeepEqual(got, want) {
		t.Errorf("addrs: %v != %v", got, want)
	}
}

type Handler struct {
	count int
	addrs []string
}

func (h *Handler) Refresh(addrs []string) {
	h.count++
	h.addrs = addrs
}

func TestServiceRefresh1(t *testing.T) {
	var h1, h2 Handler
	s := newService("name", h1.Refresh, h2.Refresh)
	kvs := map[string][]byte{
		"a1": []byte("1"),
		"a3": []byte("1"),
		"a2": []byte("1"),
	}
	addrs := []string{"a1", "a2", "a3"}
	s.Refresh(kvs)

	if got, want := h1.count, 1; got != want {
		t.Errorf("h1 count: %d != %d", got, want)
	}
	if got, want := h1.addrs, addrs; !reflect.DeepEqual(got, want) {
		t.Errorf("h1 addrs: %v != %v", got, want)
	}
	if got, want := h2.count, 1; got != want {
		t.Errorf("h2 count: %d != %d", got, want)
	}
	if got, want := h2.addrs, addrs; !reflect.DeepEqual(got, want) {
		t.Errorf("h2 addrs: %v != %v", got, want)
	}
}
