package main

import (
	"encoding/json"
	"testing"
)

type TJSONString struct {
	N int `json:",string"`
}

func TestJSONString(t *testing.T) {
	a := TJSONString{N: 10}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	t.Logf("data: %s\n", data)

	var b TJSONString
	str := `{"N":"20"}`
	if err = json.Unmarshal([]byte(str), &b); err != nil {
		t.Fatalf("unmarshal b: %v", err)
	}
	t.Logf("b: %v\n", b)

	var c TJSONString
	str = `{"N":20}`
	if err = json.Unmarshal([]byte(str), &c); err != nil {
		t.Fatalf("unmarshal c: %v", err)
	}
	t.Logf("c: %v\n", c)
}
