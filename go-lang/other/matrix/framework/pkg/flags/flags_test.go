package flags

import (
	"flag"
	"os"
	"testing"
)

func ExampleUsage1() {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.SetOutput(os.Stdout)

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
	v.Setup(f)
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

func ExampleUsage2() {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.SetOutput(os.Stdout)

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
	Setup(f, &v, "", "")
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

func TestSetup0(t *testing.T) {
	var v values
	var err error
	f := flag.NewFlagSet("", flag.ContinueOnError)

	if err = v.Setup(f); err != nil {
		t.Fatalf("setup: %v", err)
	}

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

	if err = v.Assert(); err != nil {
		t.Error(err)
	}
}

func TestSetup1(t *testing.T) {
	var v values
	var err error
	f := flag.NewFlagSet("", flag.ContinueOnError)

	if err = Setup(f, &v, "", ""); err != nil {
		t.Fatalf("setup: %v", err)
	}

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

	if err = v.Assert(); err != nil {
		t.Error(err)
	}
}

func TestSetup2(t *testing.T) {
	type V struct {
		V1 values `json:"v1" usage:"v1"`
		V2 values `json:"v2" usage:"v2"`
	}

	var v V
	var err error
	f := flag.NewFlagSet("", flag.ContinueOnError)

	if err = Setup(f, &v, "", ""); err != nil {
		t.Fatalf("setup: %v", err)
	}

	args := []string{
		"-v1.bool", "true",
		"-v1.int", "-1",
		"-v1.int8", "-8",
		"-v1.int16", "-16",
		"-v1.int32", "-32",
		"-v1.int64", "-64",
		"-v1.uint", "1",
		"-v1.uint8", "8",
		"-v1.uint16", "16",
		"-v1.uint32", "32",
		"-v1.uint64", "64",
		"-v1.uintptr", "1",
		"-v1.float32", "32.1",
		"-v1.float64", "64.1",
		"-v1.string", "1",
	}
	f.Parse(args)

	if err = v.V1.Assert(); err != nil {
		t.Error(err)
	}
}
