package iciba

import (
	"encoding/json"
	"fmt"
	"testing"
)

func ResultString(r Result) string {
	data, _ := json.MarshalIndent(r, "", "\t")
	return string(data)
}

func TestLookup(t *testing.T) {
	words := []string{"hi", "hello", "apple", "banana", "cat", "dog", "kitten", "puppy"}
	for _, w := range words {
		r, err := Lookup(w)
		if err != nil {
			t.Errorf("Lookup(%q): %v", w, err)
			continue
		}
		fmt.Println(ResultString(r))
	}
}

func BenchmarkLookup(b *testing.B) {
	words := []string{"hi", "hello", "apple", "banana", "cat", "dog", "kitten", "puppy"}
	for i := 0; i < b.N; i++ {
		Lookup(words[i%len(words)])
	}
}
