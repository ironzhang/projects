package dating

import (
	"fmt"
	"io/ioutil"
	"testing"
	"text/tabwriter"

	"github.com/ironzhang/ai-practice/knn"
)

func TestParseSampleFromCSV(t *testing.T) {
	data, err := ParseSampleFromCSV("./testdata/datingTestSet.txt")
	if err != nil {
		t.Fatal(err)
	}
	Norm(data)

	w := new(tabwriter.Writer)
	//w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	w.Init(ioutil.Discard, 0, 8, 0, '\t', 0)
	for _, s := range data {
		fmt.Fprintf(w, "%s", s.Label)
		for _, f := range s.Features {
			fmt.Fprintf(w, "\t%0.4f", f)
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func TestDating(t *testing.T) {
	data, err := ParseSampleFromCSV("./testdata/datingTestSet.txt")
	if err != nil {
		t.Fatal(err)
	}
	Norm(data)

	n := 500
	trainingData, testData := data[:n], data[n:]

	var c knn.Classifier
	c.Fit(trainingData)
	c.Test(testData, 3)
}
