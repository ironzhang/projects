package main

import (
	"os"
	"testing"
)

func TestStat(t *testing.T) {
	fi, err := os.Stat("./testdata/access.log")
	if err != nil {
		t.Fatalf("os lstat: %v", err)
	}
	t.Logf("name: %q", fi.Name())
}

func TestLstat(t *testing.T) {
	fi, err := os.Lstat("./testdata/access.log")
	if err != nil {
		t.Fatalf("os lstat: %v", err)
	}
	t.Logf("name: %q", fi.Name())
}

func TestReadlink(t *testing.T) {
	name, err := os.Readlink("./testdata/access.log")
	if err != nil {
		t.Fatalf("os readlink: %v", err)
	}
	t.Logf("name: %q", name)
}
