package model

import (
	"reflect"
	"testing"
)

func TestParseTag(t *testing.T) {
	tests := []struct {
		tag       string
		name      string
		writeable bool
	}{
		{"", "", false},
		{"name", "name", false},
		{"name,writeable", "name", true},
		{",writeable", "", true},
		{"name,readonly", "name", false},
	}
	for _, tt := range tests {
		name, writeable := parseTag(tt.tag)
		if got, want := name, tt.name; got != want {
			t.Errorf("parse %q tag: name: got(%v) != want(%v)", tt.tag, got, want)
		}
		if got, want := writeable, tt.writeable; got != want {
			t.Errorf("parse %q tag: readonly: got(%v) != want(%v)", tt.tag, got, want)
		}
	}
}

func TestTypeFields(t *testing.T) {
	type T1 struct {
		A int
		B string `json:",writeable"`
		c string
	}
	type T2 struct {
		A int    `json:"a"`
		B string `json:"b,writeable"`
		C string `json:"-"`
	}

	var t1 T1
	var t2 T2
	tests := []struct {
		t reflect.Type
		m map[string]field
	}{
		{
			t: reflect.TypeOf(t1),
			m: map[string]field{
				"A": field{"A", 0, reflect.TypeOf(t1.A), false},
				"B": field{"B", 1, reflect.TypeOf(t1.B), true},
			},
		},
		{
			t: reflect.TypeOf(t2),
			m: map[string]field{
				"a": field{"a", 0, reflect.TypeOf(t2.A), false},
				"b": field{"b", 1, reflect.TypeOf(t2.B), true},
			},
		},
	}
	for i, tt := range tests {
		m := typeFields(tt.t)
		if got, want := m, tt.m; !reflect.DeepEqual(got, want) {
			t.Errorf("tests[%d]: type fields: got(%v) != want(%v)", i, got, want)
		}
	}
}
