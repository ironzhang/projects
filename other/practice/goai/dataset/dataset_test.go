package dataset

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"
)

func TestExamplePrint(t *testing.T) {
	e := Example{
		Label:    "1",
		Instance: Instance{1, 1, 1},
	}
	e.Print(os.Stdout)
}

func TestDataSetPrint(t *testing.T) {
	d := DataSet{
		{
			Label:    "1",
			Instance: Instance{1, 1, 1},
		},
		{
			Label:    "1",
			Instance: Instance{2, 1, 1},
		},
		{
			Label:    "0",
			Instance: Instance{1, 2, 2},
		},
		{
			Label:    "0",
			Instance: Instance{2, 3, 3},
		},
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "色泽", "根蒂", "敲声", "标签")
	d.Print(w)
	w.Flush()
}

func TestDataSetValid(t *testing.T) {
	tests := []struct {
		dataset DataSet
		vaild   bool
	}{
		{
			dataset: DataSet{
				{
					Label:    "1",
					Instance: Instance{1, 1, 1},
				},
				{
					Label:    "1",
					Instance: Instance{2, 1, 1},
				},
				{
					Label:    "0",
					Instance: Instance{1, 2, 2},
				},
				{
					Label:    "0",
					Instance: Instance{2, 3, 3},
				},
			},
			vaild: true,
		},
		{
			dataset: DataSet{
				{
					Label:    "1",
					Instance: Instance{1, 1, 1},
				},
				{
					Label:    "1",
					Instance: Instance{2, 1, 1},
				},
				{
					Label:    "0",
					Instance: Instance{1, 2, 2, 3},
				},
				{
					Label:    "0",
					Instance: Instance{2, 3, 3},
				},
			},
			vaild: false,
		},
	}
	for i, tt := range tests {
		if got, want := tt.dataset.Valid(), tt.vaild; got != want {
			t.Errorf("case[%d]: %v != %v", i, got, want)
		}
	}
}
