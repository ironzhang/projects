package storage

import (
	"database/sql"
	"testing"

	"github.com/ironzhang/wordbook/cores/types"
	"github.com/ironzhang/x/dbutil"
)

func OpenTestDatabase(t *testing.T) *sql.DB {
	db, err := dbutil.OpenDatabase("mysql", DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	dbutil.TruncateTable(db, WordTableName)
	dbutil.TruncateTable(db, PriorityTableName)
	return db
}

var TestWords = []types.Word{
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
}

func SetTestWordTable(t *testing.T, tb *WordTable) {
	for i, w := range TestWords {
		if err := tb.Set(w); err != nil {
			t.Fatalf("%d: Set(%q): %v", i, w.Word, err)
		}
	}
}

var TestPriorities = []struct {
	word     string
	priority int
}{
	{word: "apple", priority: 5},
	{word: "banana", priority: 10},
}

func SetTestPriorities(t *testing.T, tb *PriorityTable) {
	for i, tt := range TestPriorities {
		if err := tb.Set(tt.word, tt.priority); err != nil {
			t.Fatalf("%d: Set(%q, %d): %v", i, tt.word, tt.priority, err)
		}
	}
}
