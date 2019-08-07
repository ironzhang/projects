package digits

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ironzhang/goai/knn"
)

func TestLoadDigitInstance(t *testing.T) {
	instance, err := LoadDigitInstance("./testdata/trainingDigits/0_0.txt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprintln(ioutil.Discard, instance)
}

func TestLoadDigitExample(t *testing.T) {
	example, err := LoadDigitExample("./testdata/trainingDigits/1_100.txt")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := example.Label, "1"; got != want {
		t.Errorf("%v != %v", got, want)
	}
}

func TestLoadDigitDataSet(t *testing.T) {
	ds, err := LoadDigitDataSet("./testdata/trainingDigits")
	if err != nil {
		t.Fatal(err)
	}
	if !ds.Valid() {
		t.Error("invalid dataset")
	}
}

func TestDigits(t *testing.T) {
	trainingData, err := LoadDigitDataSet("./testdata/trainingDigits")
	if err != nil {
		t.Fatal(err)
	}
	testData, err := LoadDigitDataSet("./testdata/testDigits")
	if err != nil {
		t.Fatal(err)
	}

	var c knn.Classifier
	c.Fit(trainingData)

	for i, example := range testData {
		got, want := c.Predict(example.Instance, 3), example.Label
		if got != want {
			fmt.Printf("case[%d]: got:%s, want:%s\n", i, got, want)
		}
	}
}
