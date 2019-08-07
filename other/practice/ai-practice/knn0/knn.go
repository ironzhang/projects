package knn

import (
	"fmt"
	"math"

	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

type Classifier struct {
	labels       []string
	trainingData *mat64.Dense
}

func (c *Classifier) Fit(labels []string, data *mat64.Dense) {
	rows, _ := data.Dims()
	if rows != len(labels) {
		panic(matrix.ErrShape)
	}

	c.labels = labels
	c.trainingData = data
}

func (c *Classifier) Predict(v *mat64.Vector, k int) string {
	rows := len(c.labels)
	distances := make(map[int]float64)
	for i := 0; i < rows; i++ {
		r := c.trainingData.RowView(i)
		distances[i] = distance(r, v)
	}
	return c.vote(rank(distances, k))
}

func (c *Classifier) PrecisionTest(labels []string, vectors []*mat64.Vector, k int) {
	if len(labels) != len(vectors) {
		panic(matrix.ErrShape)
	}
	var errcnt float64
	for i, v := range vectors {
		got := c.Predict(v, k)
		if got != labels[i] {
			errcnt += 1
		}
	}
	fmt.Printf("Precision: %f\n", 1-errcnt/float64(len(labels)))
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
