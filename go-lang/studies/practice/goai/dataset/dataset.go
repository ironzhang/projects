package dataset

import (
	"fmt"
	"io"
)

// Instance 示例
type Instance []float64

func (i Instance) Len() int {
	return len(i)
}

func (i Instance) Slice() []float64 {
	return []float64(i)
}

// Example 样例
type Example struct {
	Label    string
	Instance Instance
}

func (e *Example) Print(w io.Writer) {
	for _, v := range e.Instance {
		fmt.Fprintf(w, "%f\t", v)
	}
	fmt.Fprintf(w, "%s\n", e.Label)
}

// DataSet 数据集
type DataSet []Example

func (s DataSet) Print(w io.Writer) {
	for _, e := range s {
		e.Print(w)
	}
}

func (s DataSet) Valid() bool {
	n := len(s)
	if n <= 0 {
		return false
	}
	v := len(s[0].Instance)
	for i := 1; i < n; i++ {
		if len(s[i].Instance) != v {
			return false
		}
	}
	return true
}

func (s DataSet) Rows() int {
	return len(s)
}

func (s DataSet) Cols() int {
	return len(s[0].Instance)
}
