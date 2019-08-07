package model

import (
	"encoding"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ironzhang/matrix/jsoncfg"
)

func TestIndirect(t *testing.T) {
	var i int
	var s string
	var a interface{}
	var b interface{}
	var c io.Writer
	var d io.Writer
	var e interface{}

	b = &i
	d = os.Stdout
	e = os.Stdout

	tests := []struct {
		v reflect.Value
		p reflect.Value
	}{
		{reflect.ValueOf(&i), reflect.ValueOf(&i).Elem()},
		{reflect.ValueOf(&s), reflect.ValueOf(&s).Elem()},
		{reflect.ValueOf(&a), reflect.ValueOf(&a).Elem()},
		{reflect.ValueOf(&b), reflect.ValueOf(&i).Elem()},
		{reflect.ValueOf(b), reflect.ValueOf(&i).Elem()},
		{reflect.ValueOf(&c), reflect.ValueOf(&c).Elem()},
		{reflect.ValueOf(&d), reflect.ValueOf(os.Stdout).Elem()},
		{reflect.ValueOf(&e), reflect.ValueOf(os.Stdout).Elem()},
		{reflect.ValueOf(e), reflect.ValueOf(os.Stdout).Elem()},
	}
	for i, tt := range tests {
		p := indirect(tt.v)
		if got, want := p, tt.p; got != want {
			t.Errorf("tests[%d]: got(%v) != want(%v)", i, got, want)
		}
	}
}

func TestIndirectUnmarshaler(t *testing.T) {
	var i int
	var s string
	var a interface{}
	var b interface{}
	var d jsoncfg.Duration

	b = &i

	tests := []struct {
		v  interface{}
		u  json.Unmarshaler
		tu encoding.TextUnmarshaler
		rv reflect.Value
	}{
		{&i, nil, nil, reflect.ValueOf(&i).Elem()},
		{&s, nil, nil, reflect.ValueOf(&s).Elem()},
		{&a, nil, nil, reflect.ValueOf(&a).Elem()},
		{&b, nil, nil, reflect.ValueOf(&i).Elem()},
		{b, nil, nil, reflect.ValueOf(&i).Elem()},
		{&d, nil, &d, reflect.Value{}},
	}
	for i, tt := range tests {
		u, tu, rv := indirectUnmarshaler(reflect.ValueOf(tt.v))
		if got, want := u, tt.u; got != want {
			t.Errorf("tests[%d]: Unmarshaler: got(%v) != want(%v)", i, got, want)
		}
		if got, want := tu, tt.tu; got != want {
			t.Errorf("tests[%d]: TextUnmarshaler: got(%v) != want(%v)", i, got, want)
		}
		if got, want := rv, tt.rv; got != want {
			t.Errorf("tests[%d]: Value: got(%v) != want(%v)", i, got, want)
		}
	}
}

func TestSetValue(t *testing.T) {
	if true {
		casename := "SetHelloString"

		var x string
		var y string = "hello"

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, y; got != want {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetDurationString"

		var x jsoncfg.Duration
		var y string = "10m30s"
		var w = jsoncfg.Duration(10*time.Minute + 30*time.Second)

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; got != want {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetArrayToArray"

		var x = [5]int{5}
		var y = [10]int{1, 2, 3}
		var w = [5]int{1, 2, 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; got != want {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetArrayToSlice"

		var x = []int{5, 4, 3, 2, 1}
		var y = [3]int{1, 2, 3}
		var w = []int{1, 2, 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetSliceToArray"

		var x = [5]int{5, 4, 3, 2, 1}
		var y = []int{1, 2, 3}
		var w = [5]int{1, 2, 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; got != want {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetSliceToSlice"

		var x = []int{5, 4, 3, 2, 1}
		var y = []int{1, 2, 3}
		var w = []int{1, 2, 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetSliceToInterface"

		var x interface{}
		var y = []int{1, 2, 3}
		var w = []int{1, 2, 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetMapToInterface"

		var x interface{}
		var y = map[string]int{"1": 1, "2": 2}
		var w = y

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetMapToMap"

		var x map[string]int
		var y = map[string]int{"1": 1, "2": 2}
		var w = y

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetMapToStruct1"

		type T struct {
			A string `json:",writeable"`
			B int    `json:",writeable"`
		}

		var x T
		var y = map[string]interface{}{"1": 1, "2": 2, "A": "A", "B": 3}
		var w = T{A: "A", B: 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetMapToStruct2"

		type T struct {
			A map[string]int    `json:",writeable"`
			B map[string]string `json:",writeable"`
			C map[int]int       `json:",writeable"`
		}

		var x T
		var y = map[string]interface{}{
			"A": map[string]int{
				"A1": 1,
				"A2": 2,
			},
			"B": map[string]string{
				"B1": "1",
				"B2": "2",
			},
			"C": map[int]int{
				1: 1,
				2: 2,
			},
		}
		var w = T{
			A: map[string]int{
				"A1": 1,
				"A2": 2,
			},
			B: map[string]string{
				"B1": "1",
				"B2": "2",
			},
			C: map[int]int{
				1: 1,
				2: 2,
			},
		}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetMapToStruct3"

		type T struct {
			A string  `json:",writeable"`
			B int     `json:",writeable"`
			C float64 `json:"C,writeable"`
			D string  `json:",readonly"`
			E string  `json:"-"`
			f string
		}

		var x = T{D: "readonly", E: "nochange", f: "nochange"}
		var y = map[string]interface{}{
			"A": "A",
			"B": 3,
			"C": 4.1,
			"D": "D",
			"E": "E",
			"f": "f",
		}
		var w = T{
			A: "A",
			B: 3,
			C: 4.1,
			D: "readonly",
			E: "nochange",
			f: "nochange",
		}
		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetStructToInterface"

		type T struct {
			A string
			B int
		}

		var x interface{}
		var y = T{A: "A", B: 3}
		var w = T{A: "A", B: 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetStructToMap"

		type T struct {
			A string
			B int
		}

		var x map[string]interface{}
		var y = T{A: "A", B: 3}
		var w = map[string]interface{}{"A": "A", "B": 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}

	if true {
		casename := "SetStructToStruct"

		type T1 struct {
			A string `json:",writeable"`
			B int    `json:",writeable"`
		}

		type T2 struct {
			X string `json:"A"`
			Y int    `json:"B"`
		}

		var x T1
		var y = T2{X: "A", Y: 3}
		var w = T1{A: "A", B: 3}

		if err := setValue(&x, y); err != nil {
			t.Fatalf("%s: %v", casename, err)
		}
		if got, want := x, w; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %v != %v", casename, got, want)
		}
	}
}
