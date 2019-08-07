package dating

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/ironzhang/ai-practice/knn"
)

func ParseSampleFromCSV(filename string) ([]knn.Sample, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = '\t'

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	samples := make([]knn.Sample, len(records))
	for i, record := range records {
		for j, field := range record {
			if j+1 == len(record) {
				samples[i].Label = field
			} else {
				f64, err := strconv.ParseFloat(field, 64)
				if err != nil {
					return nil, err
				}
				samples[i].Features = append(samples[i].Features, f64)
			}
		}
	}
	return samples, nil
}

func Norm(samples []knn.Sample) {
	if len(samples) <= 0 {
		return
	}

	rows := len(samples)
	cols := len(samples[0].Features)
	for c := 0; c < cols; c++ {
		max := samples[0].Features[c]
		min := samples[0].Features[c]
		for r := 1; r < rows; r++ {
			if max < samples[r].Features[c] {
				max = samples[r].Features[c]
			}
			if min > samples[r].Features[c] {
				min = samples[r].Features[c]
			}
		}
		for r := 0; r < rows; r++ {
			samples[r].Features[c] = (samples[r].Features[c] - min) / (max - min)
		}
	}
}
