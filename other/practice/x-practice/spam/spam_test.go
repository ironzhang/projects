package spam

import (
	"fmt"
	"testing"

	"k8s.io/kubernetes/pkg/util/rand"

	"github.com/ironzhang/golang/jsoncfg"
)

var (
	testDataFile = "./testdata/testdata.json"
)

type TestData struct {
	Message  string
	Patterns []string
}

func NewRandomPatterns(n int) []string {
	patterns := make([]string, 0, n)
	for i := 0; i < n; i++ {
		patterns = append(patterns, rand.String(50))
	}
	return patterns
}

func LoadTestData() (TestData, error) {
	var td TestData
	err := jsoncfg.LoadFromFile(testDataFile, &td)
	return td, err
}

func TestProduceTestData(t *testing.T) {
	if false {
		td := TestData{
			Message:  rand.String(4000),
			Patterns: NewRandomPatterns(1000),
		}
		err := jsoncfg.WriteToFile(testDataFile, td)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("produce test data success")
	}
}

func TestEngineIsSpam(t *testing.T) {
	patterns := []string{
		"Make millions in your spare time",
		"AAA-rated",
		"BBB-rated",
		"CCC-rated",
		"DDD-rated",
		"EEE-rated",
		"FFF-rated",
		"GGG-rated",
		"XXX-rated",
	}

	tests := []struct {
		message string
		want    bool
	}{
		{"", false},
		{"hello", false},
		{"hello AAA-rated", true},
		{"你好", false},
		{"你好 Make millions in your spare time", true},
	}

	e := NewEngine(patterns)
	for _, test := range tests {
		got := e.IsSpam(test.message)
		if got != test.want {
			t.Errorf("check %q message is spam, got:%t, want:%t", test.message, got, test.want)
		}
	}
}

func BenchmarkTest1(b *testing.B) {
	var (
		patterns = []string{
			"Make millions in your spare time",
			"AAA-rated",
			"BBB-rated",
			"CCC-rated",
			"DDD-rated",
			"EEE-rated",
			"FFF-rated",
			"GGG-rated",
			"XXX-rated",
		}

		message = `When learning a new language, there are three things that you need to understand.
The first and most important is the abstract model that the language presents.
The next is the concrete syntax. Finally, you need to learn your way around the standard libraries and the common idioms of the language.`
	)

	e := NewEngine(patterns)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.IsSpam(message)
	}
}

func BenchmarkTest2(b *testing.B) {
	td, err := LoadTestData()
	if err != nil {
		b.Fatal(err)
	}

	e := NewEngine(td.Patterns)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.IsSpam(td.Message)
	}
}
