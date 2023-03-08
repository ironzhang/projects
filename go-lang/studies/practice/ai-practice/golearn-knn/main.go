package main

import (
	"fmt"
	"log"

	"github.com/gonum/matrix/mat64"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

func main2() {
	rawData, err := base.ParseCSVToInstances("./testdata/testdata.csv", false)
	if err != nil {
		log.Fatalf("parse training csv to instances: %v", err)
		return
	}
	trainingData, testData := base.InstancesTrainTestSplit(rawData, 0.5)

	model := knn.NewKnnClassifier("euclidean", "linear", 2)
	model.Fit(trainingData)
	predictions, err := model.Predict(testData)
	if err != nil {
		log.Fatalf("model predict: %v", err)
		return
	}
	fmt.Println(predictions)

	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		log.Fatalf("Unable to get confusion matrix: %v", err)
		return
	}
	fmt.Println(evaluation.GetSummary(confusionMat))
}

func newTrainingData() (values []float64, numbers []float64, rows int, cols int) {
	values = []float64{1, 1, 2, 2}
	numbers = []float64{
		1, 1.1,
		1, 1,
		0, 0,
		0, 0.1,
	}
	rows = 4
	cols = 2
	return
}

func main() {
	knnReg := knn.NewKnnRegressor("euclidean")
	knnReg.Fit(newTrainingData())
	vector := mat64.NewDense(1, 2, []float64{0, 0})
	value := knnReg.Predict(vector, 3)
	fmt.Printf("%f\n", value)
}
