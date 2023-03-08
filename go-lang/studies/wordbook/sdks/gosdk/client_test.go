package gosdk

import "testing"

func TestClientListWords(t *testing.T) {
	c := NewClient("http://localhost:8080")
	words, err := c.ListWords(0, -1)
	if err != nil {
		t.Fatalf("ListWords: %v", err)
	}
	t.Logf("words: %v", words)
}

func TestClientLookupWord(t *testing.T) {
	c := NewClient("http://localhost:8080")
	w, err := c.LookupWord("hi")
	if err != nil {
		t.Fatalf("LookupWord: %v", err)
	}
	t.Logf("word: %v", w)
}

func TestClientAdjustWordPriority(t *testing.T) {
	c := NewClient("http://localhost:8080")
	if err := c.AdjustWordPriority("hi", -1); err != nil {
		t.Fatalf("AdjustWordPriority: %v", err)
	}
}
