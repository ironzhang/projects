package storage

import (
	"testing"

	"github.com/ironzhang/x/dbutil"
)

func TestPriorityTableSet(t *testing.T) {
	db := OpenTestDatabase(t)
	tb := &PriorityTable{db: db}

	tests := []struct {
		word     string
		priority int
	}{
		{
			word:     "apple",
			priority: 5,
		},
		{
			word:     "apple",
			priority: 4,
		},
		{
			word:     "banana",
			priority: 10,
		},
		{
			word:     "banana",
			priority: 20,
		},
	}
	for i, tt := range tests {
		if err := tb.Set(tt.word, tt.priority); err != nil {
			t.Fatalf("%d: SetPriority(%q, %d): %v", i, tt.word, tt.priority, err)
		}
	}

	count, err := dbutil.Count(db, PriorityTableName)
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if got, want := count, 2; got != want {
		t.Errorf("Count: got %v, want %v", got, want)
	}
}

func TestPriorityTableGet(t *testing.T) {
	db := OpenTestDatabase(t)
	tb := &PriorityTable{db: db}
	SetTestPriorities(t, tb)

	tests := []struct {
		word     string
		ok       bool
		priority int
	}{
		{
			word:     "apple",
			ok:       true,
			priority: 5,
		},
		{
			word:     "banana",
			ok:       true,
			priority: 10,
		},
		{
			word: "apples",
			ok:   false,
		},
	}
	for i, tt := range tests {
		priority, ok, err := tb.Get(tt.word)
		if err != nil {
			t.Fatalf("%d: Get(%q): %v", i, tt.word, err)
		}
		if got, want := ok, tt.ok; got != want {
			t.Fatalf("%d: Get(%q): ok: got %v, want %v", i, tt.word, got, want)
		}
		if ok {
			if got, want := priority, tt.priority; got != want {
				t.Fatalf("%d: Get(%q): priority: got %v, want %v", i, tt.word, got, want)
			}
		}
	}
}
