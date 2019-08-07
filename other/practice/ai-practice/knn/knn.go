package knn

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

type Sample struct {
	Label    string
	Features []float64
}

type Classifier struct {
	labels       []string
	trainingData *mat64.Dense
}

func (c *Classifier) Fit(samples []Sample) {
	if len(samples) <= 0 {
		panic(matrix.ErrZeroLength)
	}

	rows := len(samples)
	cols := len(samples[0].Features)
	labs := make([]string, rows)
	data := mat64.NewDense(rows, cols, nil)
	for r := 0; r < rows; r++ {
		s := samples[r]
		if len(s.Features) != cols {
			panic(matrix.ErrShape)
		}
		labs[r] = s.Label
		for c := 0; c < cols; c++ {
			data.Set(r, c, s.Features[c])
		}
	}
	c.labels = labs
	c.trainingData = data
}

func (c *Classifier) Predict(features []float64, k int) string {
	v := mat64.NewVector(len(features), features)
	rows := len(c.labels)
	distances := make(map[int]float64)
	for i := 0; i < rows; i++ {
		r := c.trainingData.RowView(i)
		distances[i] = distance(r, v)
	}
	return c.vote(rank(distances, k))
}

func distance(x, y *mat64.Vector) float64 {
	v := mat64.NewVector(x.Len(), nil)
	v.SubVec(x, y)
	p := innerProduct(v, v)
	return math.Sqrt(p)
}

func innerProduct(x, y *mat64.Vector) float64 {
	v := mat64.NewVector(x.Len(), nil)
	v.MulElemVec(x, y)
	return mat64.Sum(v)
}

func (c *Classifier) vote(values []int) string {
	max := ""
	counts := make(map[string]int)
	for _, row := range values {
		label := c.labels[row]
		counts[label]++
		if counts[label] > counts[max] {
			max = label
		}
	}
	return max
}

func (c *Classifier) Test(samples []Sample, k int) {
	rights := 0
	wrongs := 0
	rightm := make(map[string]int)
	wrongm := make(map[string]int)
	labels := make(map[string]struct{})
	for i, s := range samples {
		if got, want := c.Predict(s.Features, k), s.Label; got == want {
			rights++
			rightm[s.Label]++
		} else {
			wrongs++
			wrongm[s.Label]++

			fmt.Fprintf(os.Stdout, "samples[%d] predict wrong, got(%s) != want(%s)\n", i, got, want)
		}
		labels[s.Label] = struct{}{}
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintf(w, "Label\tRight\tWrong\tPrecision\n")
	fmt.Fprintf(w, "-----\t-----\t-----\t---------\n")
	for label := range labels {
		right := rightm[label]
		wrong := wrongm[label]
		total := right + wrong
		fmt.Fprintf(w, "%s\t%d\t%d\t%f\n", label, right, wrong, float64(right)/float64(total))
	}
	fmt.Fprintf(w, "-----\t-----\t-----\t---------\n")
	total := rights + wrongs
	fmt.Fprintf(w, "%s\t%d\t%d\t%f\n", "total", rights, wrongs, float64(rights)/float64(total))
	w.Flush()
}
