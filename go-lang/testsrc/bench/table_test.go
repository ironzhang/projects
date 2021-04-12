package main

import (
	"testing"
)

type HashTable struct {
	m map[int]struct{}
}

func NewHashTable(size int) *HashTable {
	return &HashTable{m: make(map[int]struct{}, size)}
}

func (t *HashTable) Add(x int) {
	t.m[x] = struct{}{}
}

func (t *HashTable) Find(x int) bool {
	_, ok := t.m[x]
	return ok
}

type ArrayTable struct {
	a []int
}

func NewArrayTable(size int) *ArrayTable {
	return &ArrayTable{a: make([]int, 0, size)}
}

func (t *ArrayTable) Add(x int) {
	for _, v := range t.a {
		if v == x {
			return
		}
	}
	t.a = append(t.a, x)
}

func (t *ArrayTable) Find(x int) bool {
	for _, v := range t.a {
		if v == x {
			return true
		}
	}
	return false
}

func NewTestHashTable(n int) *HashTable {
	t := NewHashTable(n)
	for i := 0; i < n; i++ {
		t.Add(i)
	}
	return t
}

func NewTestArrayTable(n int) *ArrayTable {
	t := NewArrayTable(n)
	for i := 0; i < n; i++ {
		t.Add(i)
	}
	return t
}

func BenchmarkHashTableFind_10(b *testing.B) {
	const N = 10
	tb := NewTestHashTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkArrayTableFind_10(b *testing.B) {
	const N = 10
	tb := NewTestArrayTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkHashTableFind_20(b *testing.B) {
	const N = 20
	tb := NewTestHashTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkArrayTableFind_20(b *testing.B) {
	const N = 20
	tb := NewTestArrayTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkHashTableFind_30(b *testing.B) {
	const N = 30
	tb := NewTestHashTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkArrayTableFind_30(b *testing.B) {
	const N = 30
	tb := NewTestArrayTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkHashTableFind_40(b *testing.B) {
	const N = 40
	tb := NewTestHashTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkArrayTableFind_40(b *testing.B) {
	const N = 40
	tb := NewTestArrayTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkHashTableFind_50(b *testing.B) {
	const N = 50
	tb := NewTestHashTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkArrayTableFind_50(b *testing.B) {
	const N = 50
	tb := NewTestArrayTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkHashTableFind_100(b *testing.B) {
	const N = 100
	tb := NewTestHashTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}

func BenchmarkArrayTableFind_100(b *testing.B) {
	const N = 100
	tb := NewTestArrayTable(N)
	for i := 0; i < b.N; i++ {
		x := i % (N * 2)
		tb.Find(x)
	}
}
