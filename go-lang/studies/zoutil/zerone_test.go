package zoutil

import (
	"testing"
)

func TestZerone(t *testing.T) {
	opts := DefaultOptions
	z, err := Open("test-0", opts)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer z.Close()

	c1, err := z.NewClient("service")
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	c2, err := z.NewClient("service")
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	if c1 != c2 {
		t.Fatalf("client: %p != %p", c1, c2)
	}

	s1, err := z.NewServer("service")
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	s2, err := z.NewServer("service")
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	if s1 != s2 {
		t.Fatalf("server: %p != %p", s1, s2)
	}
}
