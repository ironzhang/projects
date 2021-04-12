package storage

import (
	"reflect"
	"testing"

	"github.com/ironzhang/wordbook/cores/types"
	"github.com/ironzhang/x/dbutil"
)

func TestWordTableSet(t *testing.T) {
	db := OpenTestDatabase(t)
	tb := &WordTable{db: db}

	words := []types.Word{
		{
			Word: "apple",
		},
		{
			Word: "apple",
			PhEN: "[ˈæpl]",
			PhAM: "[ˈæpəl]",
		},
		{
			Word: "banana",
			PhEN: "[bəˈnɑ:nə]",
			PhAM: "[bəˈnænə]",
		},
		{
			Word: "banana",
		},
	}
	for i, w := range words {
		if err := tb.Set(w); err != nil {
			t.Fatalf("%d: Set(%q): %v", i, w.Word, err)
		}
	}

	count, err := dbutil.Count(db, WordTableName)
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if got, want := count, 2; got != want {
		t.Errorf("Count: got %v, want %v", got, want)
	}
}

func TestWordTableGet(t *testing.T) {
	db := OpenTestDatabase(t)
	tb := &WordTable{db: db}
	SetTestWordTable(t, tb)

	tests := []struct {
		w  types.Word
		ok bool
	}{
		{
			w: types.Word{
				Word: "apple",
				PhEN: "[ˈæpl]",
				PhAM: "[ˈæpəl]",
			},
			ok: true,
		},
		{
			w: types.Word{
				Word: "apples",
			},
			ok: false,
		},
		{
			w: types.Word{
				Word: "banana",
				PhEN: "[bəˈnɑ:nə]",
				PhAM: "[bəˈnænə]",
			},
			ok: true,
		},
		{
			w: types.Word{
				Word: "bananas",
			},
			ok: false,
		},
	}
	for i, tt := range tests {
		w, ok, err := tb.Get(tt.w.Word)
		if err != nil {
			t.Fatalf("%d: Get(%q): %v", i, tt.w.Word, err)
		}
		if got, want := ok, tt.ok; got != want {
			t.Fatalf("%d: Get(%q): ok: got %v, want %v", i, tt.w.Word, got, want)
		}
		if ok {
			if got, want := w, tt.w; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d: Get(%q): word: got %v, want %v", i, tt.w.Word, got, want)
			} else {
				t.Logf("%d: Get(%q): word: %v", i, tt.w.Word, got)
			}
		}
	}
}

func TestWordTableRemove(t *testing.T) {
	db := OpenTestDatabase(t)
	tb := &WordTable{db: db}
	SetTestWordTable(t, tb)

	words := []string{"apple", "banana"}
	for i, w := range words {
		if err := tb.Remove(w); err != nil {
			t.Fatalf("%d: Remove(%q): %v", i, w, err)
		}
		_, ok, err := tb.Get(w)
		if err != nil {
			t.Fatalf("%d: Get(%q): %v", i, w, err)
		}
		if got, want := ok, false; got != want {
			t.Fatalf("%d: Get(%q): got %v, want %v", i, w, got, want)
		}
	}
}

func TestWordTableList(t *testing.T) {
	db := OpenTestDatabase(t)
	tb := &WordTable{db: db}
	SetTestWordTable(t, tb)

	tests := []struct {
		offset int
		limit  int
		words  []types.Word
	}{
		{
			offset: 0,
			limit:  -1,
			words:  TestWords,
		},
		{
			offset: 0,
			limit:  1,
			words:  TestWords[:1],
		},
		{
			offset: 1,
			limit:  -1,
			words:  TestWords[1:],
		},
	}
	for i, tt := range tests {
		words, err := tb.List(tt.offset, tt.limit)
		if err != nil {
			t.Fatalf("%d: List: %v", i, err)
		}
		if got, want := words, tt.words; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: List: got %v, want %v", i, got, want)
		}
		t.Logf("%d: List: got %v", i, words)
	}
}

func TestWordTableListByPriority(t *testing.T) {
	db := OpenTestDatabase(t)
	tb := &WordTable{db: db}
	ptb := &PriorityTable{db: db}
	SetTestWordTable(t, tb)
	SetTestPriorities(t, ptb)

	tests := []struct {
		offset int
		limit  int
		words  []types.Word
	}{
		{
			offset: 0,
			limit:  -1,
			words:  []types.Word{TestWords[1], TestWords[0]},
		},
		{
			offset: 0,
			limit:  1,
			words:  []types.Word{TestWords[1]},
		},
		{
			offset: 1,
			limit:  -1,
			words:  []types.Word{TestWords[0]},
		},
	}
	for i, tt := range tests {
		words, err := tb.ListByPriority(tt.offset, tt.limit)
		if err != nil {
			t.Fatalf("%d: ListByPriority: %v", i, err)
		}
		if got, want := words, tt.words; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: ListByPriority: got %v, want %v", i, got, want)
		}
		t.Logf("%d: ListByPriority: got %v", i, words)
	}
}
