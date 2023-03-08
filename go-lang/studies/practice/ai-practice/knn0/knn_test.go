package knn

import (
	"fmt"
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestClassifier(t *testing.T) {
	tests := []struct {
		labels       []string
		trainingData *mat64.Dense
		testData     *mat64.Vector
		k            int
	}{
		{
			labels: []string{"A", "A", "B", "B"},
			trainingData: mat64.NewDense(4, 2, []float64{
				1, 1.1,
				1, 1,
				0, 0,
				0, 0.1,
			}),
			testData: mat64.NewVector(2, []float64{0, 0}),
			k:        3,
		},

		{
			labels: []string{"爱情片", "爱情片", "爱情片", "动作片", "动作片", "动作片"},
			trainingData: mat64.NewDense(6, 2, []float64{
				3, 104,
				2, 100,
				1, 81,
				101, 10,
				99, 5,
				98, 2,
			}),
			testData: mat64.NewVector(2, []float64{18, 90}),
			k:        3,
		},
	}
	for i, tt := range tests {
		var c Classifier
		c.Fit(tt.labels, tt.trainingData)
		predictions := c.Predict(tt.testData, tt.k)
		fmt.Printf("case%d: %s\n", i, predictions)
	}
}

func TestClassifierPrecision(t *testing.T) {
	tests := []struct {
		labels       []string
		trainingData *mat64.Dense
		testData     []*mat64.Vector
		testLabels   []string
		k            int
	}{
		{
			labels: []string{"A", "A", "B", "B"},
			trainingData: mat64.NewDense(4, 2, []float64{
				1, 1.1,
				1, 1,
				0, 0,
				0, 0.1,
			}),
			testData: []*mat64.Vector{
				mat64.NewVector(2, []float64{0, 0}),
			},
			testLabels: []string{"B"},
			k:          3,
		},

		{
			labels: []string{"爱情片", "爱情片", "爱情片", "动作片", "动作片", "动作片"},
			trainingData: mat64.NewDense(6, 2, []float64{
				3, 104,
				2, 100,
				1, 81,
				101, 10,
				99, 5,
				98, 2,
			}),
			testData: []*mat64.Vector{
				mat64.NewVector(2, []float64{18, 90}),
			},
			testLabels: []string{"爱情片"},
			k:          3,
		},
	}
	for _, tt := range tests {
		var c Classifier
		c.Fit(tt.labels, tt.trainingData)
		c.PrecisionTest(tt.testLabels, tt.testData, tt.k)
	}
}
