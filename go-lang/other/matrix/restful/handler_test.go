package restful

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type AStruct struct{}
type bStruct struct{}
type _CStruct struct{}

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

func TestIsExported(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", false},
		{"a", false},
		{"ab", false},
		{"A", true},
		{"Ab", true},
		{"_a", false},
		{"A_", true},
		{"中文", false},
	}
	for _, tt := range tests {
		if got, want := isExported(tt.name), tt.want; got != want {
			t.Errorf("name:%q, got(%t) != want(%t)", tt.name, got, want)
		}
	}
}

func TestIsExportedOrBuiltinType(t *testing.T) {
	var a AStruct
	var b bStruct
	var c _CStruct
	var i int
	var s string
	var v reflect.Value

	tests := []struct {
		value interface{}
		want  bool
	}{
		{a, true},
		{&a, true},
		{b, false},
		{&b, false},
		{c, false},
		{&c, false},
		{i, true},
		{&i, true},
		{s, true},
		{&s, true},
		{v, true},
		{&v, true},
	}
	for _, tt := range tests {
		if got, want := isExportedOrBuiltinType(reflect.TypeOf(tt.value)), tt.want; got != want {
			t.Errorf("type:%T, got(%t) != want(%t)", tt.value, got, want)
		}
	}
}

func TestCheckIns(t *testing.T) {
	wrongs := []interface{}{
		func() {},
		func(interface{}) {},
		func(interface{}, interface{}) {},
		func(interface{}, interface{}, interface{}) {},
		func(interface{}, interface{}, interface{}, interface{}) {},
		func(Context, interface{}, interface{}, interface{}) {},
		func(Context, url.Values, bStruct, interface{}) {},
		func(Context, url.Values, interface{}, int) {},
		func(Context, url.Values, interface{}, AStruct) {},
		func(Context, url.Values, interface{}, *bStruct) {},
	}
	for i, f := range wrongs {
		if _, _, _, _, err := checkIns(reflect.TypeOf(f)); err == nil {
			t.Errorf("wrongs[%d]: checkIns %T, expect error but not", i, f)
		} else {
			t.Logf("wrongs[%d]: checkIns: %v", i, err)
		}
	}

	rights := []interface{}{
		func(context.Context, url.Values, interface{}, interface{}) {},
		func(Context, url.Values, interface{}, interface{}) {},
		func(Context, url.Values, int, interface{}) {},
		func(Context, url.Values, *int, interface{}) {},
		func(Context, url.Values, AStruct, interface{}) {},
		func(Context, url.Values, *AStruct, interface{}) {},
		func(Context, url.Values, interface{}, *int) {},
		func(Context, url.Values, interface{}, *AStruct) {},
	}
	for i, f := range rights {
		in0, in1, in2, in3, err := checkIns(reflect.TypeOf(f))
		if err != nil {
			t.Errorf("rights[%d]: checkIns: %v", i, err)
		}
		if in0 != reflect.TypeOf(f).In(0) {
			t.Errorf("rights[%d]: in0: %s != %s", i, in0, reflect.TypeOf(f).In(0))
		}
		if in1 != reflect.TypeOf(f).In(1) {
			t.Errorf("rights[%d]: in1: %s != %s", i, in1, reflect.TypeOf(f).In(1))
		}
		if in2 != reflect.TypeOf(f).In(2) {
			t.Errorf("rights[%d]: in2: %s != %s", i, in2, reflect.TypeOf(f).In(2))
		}
		if in3 != reflect.TypeOf(f).In(3) {
			t.Errorf("rights[%d]: in3: %s != %s", i, in2, reflect.TypeOf(f).In(2))
		}
	}
}

func TestCheckOuts(t *testing.T) {
	wrongs := []interface{}{
		func() {},
		func() int { return 0 },
		func() (error, error) { return nil, nil },
	}
	for i, f := range wrongs {
		if err := checkOuts(reflect.TypeOf(f)); err == nil {
			t.Errorf("wrongs[%d]: checkOuts %T, expect error but not", i, f)
		} else {
			t.Logf("wrongs[%d]: checkOuts: %v", i, err)
		}
	}

	rights := []interface{}{
		func() error { return nil },
	}
	for i, f := range rights {
		if err := checkOuts(reflect.TypeOf(f)); err != nil {
			t.Errorf("rights[%d]: checkOuts: %v", i, err)
		}
	}
}

func TestParseHandler(t *testing.T) {
	rights := []interface{}{
		func(context.Context, url.Values, interface{}, interface{}) error { return nil },
		func(context.Context, url.Values, AStruct, interface{}) error { return nil },
		func(Context, url.Values, *AStruct, *AStruct) error { return nil },
	}
	for i, f := range rights {
		h, err := parseHandler(f)
		if err != nil {
			t.Errorf("rights[%d]: parse handler: %v", i, err)
			continue
		}
		if h.value != reflect.ValueOf(f) {
			t.Errorf("rights[%d]: value: %v != %v", i, h.value, reflect.ValueOf(f))
			continue
		}
		if h.args != reflect.TypeOf(f).In(2) {
			t.Errorf("rights[%d]: args type: %s != %s", i, h.args, reflect.TypeOf(f).In(1))
			continue
		}
		if h.reply != reflect.TypeOf(f).In(3) {
			t.Errorf("rights[%d]: reply type: %s != %s", i, h.reply, reflect.TypeOf(f).In(2))
			continue
		}
	}
}

func TestHandleReturnNil(t *testing.T) {
	tests := []interface{}{
		func(context.Context, url.Values, interface{}, interface{}) error { return nil },
		func(Context, url.Values, interface{}, interface{}) error { return nil },
		func(Context, url.Values, int, interface{}) error { return nil },
		func(Context, url.Values, *int, interface{}) error { return nil },
		func(Context, url.Values, AStruct, interface{}) error { return nil },
		func(Context, url.Values, *AStruct, interface{}) error { return nil },
		func(Context, url.Values, interface{}, *int) error { return nil },
		func(Context, url.Values, interface{}, *AStruct) error { return nil },
	}
	for i, f := range tests {
		h, err := parseHandler(f)
		if err != nil {
			t.Errorf("tests[%d]: parse handler: %v", i, err)
			continue
		}
		err = h.Handle(context.Background(), nil, newReflectValue(h.args), newReflectValue(h.reply))
		if err != nil {
			t.Errorf("tests[%d]: handle: %v", i, err)
			continue
		}
	}
}

func TestHandleReturnErr(t *testing.T) {
	f := func(context.Context, url.Values, interface{}, interface{}) error { return io.EOF }
	h, err := parseHandler(f)
	if err != nil {
		t.Fatalf("parse handler: %v", err)
	}
	err = h.Handle(context.Background(), nil, newReflectValue(h.args), newReflectValue(h.reply))
	if err != io.EOF {
		t.Errorf("err: %v != %v", err, io.EOF)
	}
}

func TestHandler(t *testing.T) {
	type Request struct {
		A, B int
	}
	type Response struct {
		C int
	}

	var calls int
	f := func(ctx context.Context, values url.Values, req *Request, resp *Response) error {
		calls++
		resp.C = req.A + req.B
		return nil
	}
	h, err := parseHandler(f)
	if err != nil {
		t.Fatalf("parse handler: %v", err)
	}

	var req Request
	var resp Response
	req.A = 1
	req.B = 2
	if err = h.Handle(context.Background(), nil, reflect.ValueOf(&req), reflect.ValueOf(&resp)); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if calls != 1 {
		t.Errorf("calls(%d) != 1", calls)
	}
	if resp.C != req.A+req.B {
		t.Errorf("C(%d) != A(%d) + B(%d)", resp.C, req.A, req.B)
	}
}

func TestIsNilInterface(t *testing.T) {
	f := func(context.Context, url.Values, interface{}, interface{}) error { return io.EOF }
	h, err := parseHandler(f)
	if err != nil {
		t.Fatalf("parse handler: %v", err)
	}
	if got, want := isNilInterface(h.args), true; got != want {
		t.Errorf("in1: got(%t) != want(%t)", got, want)
	}
	if got, want := isNilInterface(h.reply), true; got != want {
		t.Errorf("in2: got(%t) != want(%t)", got, want)
	}

	{
		want := true
		var v interface{}
		if got := isNilInterface(reflect.TypeOf(&v).Elem()); got != want {
			t.Errorf("%T: got(%t) != want(%t)", v, got, want)
		}
	}
	{
		want := false
		var v struct{}
		if got := isNilInterface(reflect.TypeOf(&v).Elem()); got != want {
			t.Errorf("%T: got(%t) != want(%t)", v, got, want)
		}
	}
	{
		want := false
		var v int
		if got := isNilInterface(reflect.TypeOf(&v).Elem()); got != want {
			t.Errorf("%T: got(%t) != want(%t)", v, got, want)
		}
	}
	{
		want := false
		var v *int
		if got := isNilInterface(reflect.TypeOf(&v).Elem()); got != want {
			t.Errorf("%T: got(%t) != want(%t)", v, got, want)
		}
	}
	{
		want := false
		var v Context
		if got := isNilInterface(reflect.TypeOf(&v).Elem()); got != want {
			t.Errorf("%T: got(%t) != want(%t)", v, got, want)
		}
	}
}

func TestJSON(t *testing.T) {
	{
		var v interface{}
		if buf, err := json.Marshal(v); err != nil {
			t.Errorf("%T json marshal: %v", v, err)
		} else {
			t.Logf("%T buf: %s", v, buf)
		}
	}
	{
		var v struct{}
		if buf, err := json.Marshal(v); err != nil {
			t.Errorf("%T json marshal: %v", v, err)
		} else {
			t.Logf("%T buf: %s", v, buf)
		}
	}
	{
		var v int
		if buf, err := json.Marshal(v); err != nil {
			t.Errorf("%T json marshal: %v", v, err)
		} else {
			t.Logf("%T buf: %s", v, buf)
		}
	}
	{
		var v string
		if buf, err := json.Marshal(v); err != nil {
			t.Errorf("%T json marshal: %v", v, err)
		} else {
			t.Logf("%T buf: %s", v, buf)
		}
	}
	{
		var v float32
		if buf, err := json.Marshal(v); err != nil {
			t.Errorf("%T json marshal: %v", v, err)
		} else {
			t.Logf("%T buf: %s", v, buf)
		}
	}
}
