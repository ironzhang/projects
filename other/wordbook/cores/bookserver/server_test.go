package bookserver

import (
	"database/sql"
	"testing"

	"github.com/ironzhang/wordbook/cores/storage"
	"github.com/ironzhang/x/dbutil"
)

func OpenTestDatabase(t *testing.T) *sql.DB {
	db, err := dbutil.OpenDatabase("mysql", storage.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestServerLookupWord(t *testing.T) {
	db := OpenTestDatabase(t)
	s := NewServer(db)

	words := []string{
		"apple", "banana", "bananas", "cat", "duck", "elephant", "fox", "green", "hello",
		"I", "juice", "kitten", "long", "Mum", "nest",

		"raining", "play", "muddy", "puddle", "jump", "boot", "mud", "television", "inside",
		"outside", "bath", "mess", "quickly", "slowly", "clean", "garden", "love", "pig",

		"toy", "scare", "scary", "supper", "bed", "game", "throwing", "down", "draughts",
		"detective", "question", "night", "happy", "favourite", "tucked", "important",

		"friend", "sheep", "bedroom", "girl", "toy", "fairy", "princess", "wave", "magic", "wand",
		"frog", "cookie", "chocolate", "bowl", "nurse", "doctor", "listen", "cough", "heart", "loose",

		"parrot", "surprise", "Grandpa", "shy", "cake", "pet", "tea", "word", "nurse", "queen", "rose",
		"octopus", "pig", "X-ray", "snail", "van", "umbrella", "watch", "turtle", "yellow", "zebra",
		"ox", "quilt",
	}
	for i, word := range words {
		w, ok, err := s.LookupWord(word)
		if err != nil {
			t.Fatalf("%d: LookupWord(%q): %v", i, word, err)
		}
		if !ok {
			t.Errorf("%d: LookupWord(%q): not found", i, word)
		}
		t.Logf("%d: LookupWord(%q): %v", i, word, w)
	}
}
