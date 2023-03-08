package flags

import (
	"flag"
	"fmt"
	"reflect"
	"runtime"

	"github.com/ironzhang/matrix/errs"
	"github.com/ironzhang/matrix/framework/pkg/tags"
)

func Setup(f *flag.FlagSet, value interface{}, name, usage string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if e, ok := r.(error); ok {
				err = e
				return
			}
			panic(r)
		}
	}()

	flags{f}.SetupValue(name, usage, reflect.ValueOf(value).Elem())
	return
}

type flags struct {
	*flag.FlagSet
}

func (f flags) SetupValue(name, usage string, v reflect.Value) {
	switch k := v.Kind(); k {
	case reflect.Bool:
		f.Var(newBoolValue(v), name, usage)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.Var(newIntValue(v), name, usage)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		f.Var(newUintValue(v), name, usage)
	case reflect.Float32, reflect.Float64:
		f.Var(newFloatValue(v), name, usage)
	case reflect.String:
		f.Var(newStringValue(v), name, usage)
	case reflect.Struct:
		if name != "" {
			name = name + "."
		}
		if usage != "" {
			usage = usage + ": "
		}
		f.SetupStruct(name, usage, v)
	default:
		panic(errs.ErrorAt("flags.SetupValue", fmt.Errorf("unsupport %s kind", k)))
	}
}

func (f flags) SetupStruct(prefix, usage string, v reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}
		name := sf.Name
		if s := parseNameFromTag(sf.Tag); s != "" {
			name = s
		}
		if name == "-" {
			continue
		}
		f.SetupValue(prefix+name, usage+sf.Tag.Get("usage"), v.Field(i))
	}
}

func parseNameFromTag(tag reflect.StructTag) string {
	name, _ := tags.ParseTag(tag.Get("json"))
	if tags.IsValidTag(name) {
		return name
	}
	return ""
}
