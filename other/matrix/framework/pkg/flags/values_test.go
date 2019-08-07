package flags

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
)

type values struct {
	B    bool    `json:"bool" usage:"usage of bool"`
	I    int     `json:"int" usage:"usage of int"`
	I8   int8    `json:"int8" usage:"usage of int8"`
	I16  int16   `json:"int16" usage:"usage of int16"`
	I32  int32   `json:"int32" usage:"usage of int32"`
	I64  int64   `json:"int64" usage:"usage of int64"`
	U    uint    `json:"uint" usage:"usage of uint"`
	U8   uint8   `json:"uint8" usage:"usage of uint8"`
	U16  uint16  `json:"uint16" usage:"usage of uint16"`
	U32  uint32  `json:"uint32" usage:"usage of uint32"`
	U64  uint64  `json:"uint64" usage:"usage of uint64"`
	Uptr uintptr `json:"uintptr" usage:"usage of uintptr"`
	F32  float32 `json:"float32" usage:"usage of float32"`
	F64  float64 `json:"float64" usage:"usage of float64"`
	S    string  `json:"string" usage:"usage of string"`
}

func (v *values) SetupVars(f *flag.FlagSet) {
	f.Var(newBoolValue(reflect.ValueOf(&v.B).Elem()), "bool", "usage of bool")
	f.Var(newIntValue(reflect.ValueOf(&v.I).Elem()), "int", "usage of int")
	f.Var(newIntValue(reflect.ValueOf(&v.I8).Elem()), "int8", "usage of int8")
	f.Var(newIntValue(reflect.ValueOf(&v.I16).Elem()), "int16", "usage of int16")
	f.Var(newIntValue(reflect.ValueOf(&v.I32).Elem()), "int32", "usage of int32")
	f.Var(newIntValue(reflect.ValueOf(&v.I64).Elem()), "int64", "usage of int64")
	f.Var(newUintValue(reflect.ValueOf(&v.U).Elem()), "uint", "usage of uint")
	f.Var(newUintValue(reflect.ValueOf(&v.U8).Elem()), "uint8", "usage of uint8")
	f.Var(newUintValue(reflect.ValueOf(&v.U16).Elem()), "uint16", "usage of uint16")
	f.Var(newUintValue(reflect.ValueOf(&v.U32).Elem()), "uint32", "usage of uint32")
	f.Var(newUintValue(reflect.ValueOf(&v.U64).Elem()), "uint64", "usage of uint64")
	f.Var(newUintValue(reflect.ValueOf(&v.Uptr).Elem()), "uintptr", "usage of uintptr")
	f.Var(newFloatValue(reflect.ValueOf(&v.F32).Elem()), "float32", "usage of float32")
	f.Var(newFloatValue(reflect.ValueOf(&v.F64).Elem()), "float64", "usage of float64")
	f.Var(newStringValue(reflect.ValueOf(&v.S).Elem()), "string", "usage of string")
}

func (v *values) Setup(f *flag.FlagSet) (err error) {
	if err = Setup(f, &v.B, "bool", "usage of bool"); err != nil {
		return err
	}
	if err = Setup(f, &v.I, "int", "usage of int"); err != nil {
		return err
	}
	if err = Setup(f, &v.I8, "int8", "usage of int8"); err != nil {
		return err
	}
	if err = Setup(f, &v.I16, "int16", "usage of int16"); err != nil {
		return err
	}
	if err = Setup(f, &v.I32, "int32", "usage of int32"); err != nil {
		return err
	}
	if err = Setup(f, &v.I64, "int64", "usage of int64"); err != nil {
		return err
	}
	if err = Setup(f, &v.U, "uint", "usage of uint"); err != nil {
		return err
	}
	if err = Setup(f, &v.U8, "uint8", "usage of uint8"); err != nil {
		return err
	}
	if err = Setup(f, &v.U16, "uint16", "usage of uint16"); err != nil {
		return err
	}
	if err = Setup(f, &v.U32, "uint32", "usage of uint32"); err != nil {
		return err
	}
	if err = Setup(f, &v.U64, "uint64", "usage of uint64"); err != nil {
		return err
	}
	if err = Setup(f, &v.Uptr, "uintptr", "usage of uintptr"); err != nil {
		return err
	}
	if err = Setup(f, &v.F32, "float32", "usage of float32"); err != nil {
		return err
	}
	if err = Setup(f, &v.F64, "float64", "usage of float64"); err != nil {
		return err
	}
	if err = Setup(f, &v.S, "string", "usage of string"); err != nil {
		return err
	}
	return nil
}

func (v *values) Assert() error {
	if got, want := v.B, true; got != want {
		return fmt.Errorf("bool: %v != %v", got, want)
	}
	if got, want := v.I, int(-1); got != want {
		return fmt.Errorf("int: %v != %v", got, want)
	}
	if got, want := v.I8, int8(-8); got != want {
		return fmt.Errorf("int8: %v != %v", got, want)
	}
	if got, want := v.I16, int16(-16); got != want {
		return fmt.Errorf("int16: %v != %v", got, want)
	}
	if got, want := v.I32, int32(-32); got != want {
		return fmt.Errorf("int32: %v != %v", got, want)
	}
	if got, want := v.I64, int64(-64); got != want {
		return fmt.Errorf("int64: %v != %v", got, want)
	}
	if got, want := v.U, uint(1); got != want {
		return fmt.Errorf("uint: %v != %v", got, want)
	}
	if got, want := v.U8, uint8(8); got != want {
		return fmt.Errorf("uint8: %v != %v", got, want)
	}
	if got, want := v.U16, uint16(16); got != want {
		return fmt.Errorf("uint16: %v != %v", got, want)
	}
	if got, want := v.U32, uint32(32); got != want {
		return fmt.Errorf("uint32: %v != %v", got, want)
	}
	if got, want := v.U64, uint64(64); got != want {
		return fmt.Errorf("uint64: %v != %v", got, want)
	}
	if got, want := v.Uptr, uintptr(1); got != want {
		return fmt.Errorf("uintptr: %v != %v", got, want)
	}
	if got, want := v.F32, float32(32.1); got != want {
		return fmt.Errorf("float32: %v != %v", got, want)
	}
	if got, want := v.F64, float64(64.1); got != want {
		return fmt.Errorf("float64: %v != %v", got, want)
	}
	if got, want := v.S, "1"; got != want {
		return fmt.Errorf("string: %v != %v", got, want)
	}
	return nil
}

func ExampleUsage0() {
	v := values{
		B:   true,
		I:   -1,
		I8:  -8,
		I16: -16,
		I32: -32,
		I64: -64,
		U:   1,
		U8:  8,
		U16: 16,
		U32: 32,
		U64: 64,
		F32: 32.1,
		F64: 64.1,
		S:   "1",
	}
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.SetOutput(os.Stdout)
	v.SetupVars(f)
	f.Usage()

	// output:
	// Usage:
	//   -bool value
	//     	usage of bool (default true)
	//   -float32 value
	//     	usage of float32 (default 32.099998474121094)
	//   -float64 value
	//     	usage of float64 (default 64.1)
	//   -int value
	//     	usage of int (default -1)
	//   -int16 value
	//     	usage of int16 (default -16)
	//   -int32 value
	//     	usage of int32 (default -32)
	//   -int64 value
	//     	usage of int64 (default -64)
	//   -int8 value
	//     	usage of int8 (default -8)
	//   -string value
	//     	usage of string (default 1)
	//   -uint value
	//     	usage of uint (default 1)
	//   -uint16 value
	//     	usage of uint16 (default 16)
	//   -uint32 value
	//     	usage of uint32 (default 32)
	//   -uint64 value
	//     	usage of uint64 (default 64)
	//   -uint8 value
	//     	usage of uint8 (default 8)
	//   -uintptr value
	//     	usage of uintptr
}

func TestValues(t *testing.T) {
	var v values
	f := flag.NewFlagSet("", flag.ContinueOnError)
	v.SetupVars(f)

	args := []string{
		"-bool", "true",
		"-int", "-1",
		"-int8", "-8",
		"-int16", "-16",
		"-int32", "-32",
		"-int64", "-64",
		"-uint", "1",
		"-uint8", "8",
		"-uint16", "16",
		"-uint32", "32",
		"-uint64", "64",
		"-uintptr", "1",
		"-float32", "32.1",
		"-float64", "64.1",
		"-string", "1",
	}
	f.Parse(args)

	if err := v.Assert(); err != nil {
		t.Error(err)
	}
}
