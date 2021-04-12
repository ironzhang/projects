package flags

import (
	"reflect"
	"strconv"
)

type boolValue reflect.Value

func newBoolValue(v reflect.Value) boolValue {
	return boolValue(v)
}

func (v boolValue) Set(s string) error {
	x, err := strconv.ParseBool(s)
	reflect.Value(v).SetBool(x)
	return err
}

func (v boolValue) String() string {
	val := reflect.Value(v)
	if !val.IsValid() {
		return "false"
	}
	return strconv.FormatBool(val.Bool())
}

type intValue reflect.Value

func newIntValue(v reflect.Value) intValue {
	return intValue(v)
}

func (v intValue) Set(s string) error {
	x, err := strconv.ParseInt(s, 0, 64)
	reflect.Value(v).SetInt(x)
	return err
}

func (v intValue) String() string {
	val := reflect.Value(v)
	if !val.IsValid() {
		return "0"
	}
	return strconv.FormatInt(val.Int(), 10)
}

type uintValue reflect.Value

func newUintValue(v reflect.Value) uintValue {
	return uintValue(v)
}

func (v uintValue) Set(s string) error {
	x, err := strconv.ParseUint(s, 0, 64)
	reflect.Value(v).SetUint(x)
	return err
}

func (v uintValue) String() string {
	val := reflect.Value(v)
	if !val.IsValid() {
		return "0"
	}
	return strconv.FormatUint(val.Uint(), 10)
}

type floatValue reflect.Value

func newFloatValue(v reflect.Value) floatValue {
	return floatValue(v)
}

func (v floatValue) Set(s string) error {
	x, err := strconv.ParseFloat(s, 64)
	reflect.Value(v).SetFloat(x)
	return err
}

func (v floatValue) String() string {
	val := reflect.Value(v)
	if !val.IsValid() {
		return "0"
	}
	return strconv.FormatFloat(val.Float(), 'g', -1, 64)
}

type stringValue reflect.Value

func newStringValue(v reflect.Value) stringValue {
	return stringValue(v)
}

func (v stringValue) Set(s string) error {
	reflect.Value(v).SetString(s)
	return nil
}

func (v stringValue) String() string {
	val := reflect.Value(v)
	if !val.IsValid() {
		return ""
	}
	return val.String()
}
