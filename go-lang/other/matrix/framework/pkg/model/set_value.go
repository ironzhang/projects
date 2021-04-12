package model

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/ironzhang/matrix/errs"
)

var errInvalidValue = errors.New("invalid value")

func indirect(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		panic(errs.ErrorAt("indirect", errInvalidValue))
	}

	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

L:
	for {
		switch k := v.Kind(); k {
		case reflect.Ptr:
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		case reflect.Interface:
			if v.IsNil() {
				break L
			}
			v = v.Elem()
		default:
			break L
		}
	}
	return v
}

func indirectUnmarshaler(v reflect.Value) (json.Unmarshaler, encoding.TextUnmarshaler, reflect.Value) {
	if !v.IsValid() {
		panic(errs.ErrorAt("indirectUnmarshaler", errInvalidValue))
	}

	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

L:
	for {
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(json.Unmarshaler); ok {
				return u, nil, reflect.Value{}
			}
			if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
				return nil, u, reflect.Value{}
			}
		}

		switch k := v.Kind(); k {
		case reflect.Ptr:
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		case reflect.Interface:
			if v.IsNil() {
				break L
			}
			v = v.Elem()
		default:
			break L
		}
	}
	return nil, nil, v
}

type setError struct {
	src reflect.Type
	dst reflect.Type
	err error
}

func (e setError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("can not set value from %s[%s] to %s[%s]: %v", e.src.Kind(), e.src, e.dst.Kind(), e.dst, e.err)
	}
	return fmt.Sprintf("can not set value from %s[%s] to %s[%s]", e.src.Kind(), e.src, e.dst.Kind(), e.dst)
}

// SetValue set the y value to x
func setValue(x, y interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if e, ok := r.(error); ok {
				err = e
				return
			}
			if e, ok := r.(string); ok {
				err = errors.New(e)
				return
			}
			panic(r)
		}
	}()

	setState{}.set(reflect.ValueOf(x), reflect.ValueOf(y))
	return
}

type setState struct {
}

func (s setState) set(x, y reflect.Value) {
	y = indirect(y)
	switch y.Kind() {
	case reflect.Bool:
		s.setBool(x, y)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.setInt(x, y)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		s.setUint(x, y)

	case reflect.Float32, reflect.Float64:
		s.setFloat(x, y)

	case reflect.String:
		s.setString(x, y)

	case reflect.Array, reflect.Slice:
		s.setArray(x, y)

	case reflect.Map:
		s.setMap(x, y)

	case reflect.Struct:
		s.setStruct(x, y)

	default:
		s.panic("set", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) setBool(x, y reflect.Value) {
	x = indirect(x)
	switch x.Kind() {
	case reflect.Bool:
		x.SetBool(y.Bool())

	default:
		s.panic("setBool", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) setInt(x, y reflect.Value) {
	x = indirect(x)
	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x.SetInt(y.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x.SetUint(uint64(y.Int()))

	case reflect.Float32, reflect.Float64:
		x.SetFloat(float64(y.Int()))

	default:
		s.panic("setInt", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) setUint(x, y reflect.Value) {
	x = indirect(x)
	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x.SetInt(int64(y.Uint()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x.SetUint(y.Uint())

	case reflect.Float32, reflect.Float64:
		x.SetFloat(float64(y.Uint()))

	default:
		s.panic("setUint", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) setFloat(x, y reflect.Value) {
	x = indirect(x)
	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x.SetInt(int64(y.Float()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x.SetUint(uint64(y.Float()))

	case reflect.Float32, reflect.Float64:
		x.SetFloat(y.Float())

	default:
		s.panic("setFloat", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) setString(x, y reflect.Value) {
	u, ut, v := indirectUnmarshaler(x)
	if u != nil {
		b := []byte(`"` + y.String() + `"`)
		if err := u.UnmarshalJSON(b); err != nil {
			s.panic("setString", setError{src: y.Type(), dst: x.Type(), err: err})
		}
		return
	}
	if ut != nil {
		b := []byte(y.String())
		if err := ut.UnmarshalText(b); err != nil {
			s.panic("setString", setError{src: y.Type(), dst: x.Type(), err: err})
		}
		return
	}
	v.Set(y)
}

func (s setState) setArray(x, y reflect.Value) {
	x = indirect(x)
	switch x.Kind() {
	case reflect.Interface:
		if x.NumMethod() != 0 {
			s.panic("setArray", setError{src: y.Type(), dst: x.Type()})
		}
		x.Set(y)

	case reflect.Array:
		for i := 0; i < x.Len(); i++ {
			v := x.Index(i)
			v.Set(reflect.Zero(v.Type()))
		}
		for i := 0; i < x.Len() && i < y.Len(); i++ {
			s.set(x.Index(i), y.Index(i))
		}

	case reflect.Slice:
		ns := reflect.MakeSlice(x.Type(), y.Len(), y.Len())
		for i := 0; i < y.Len(); i++ {
			s.set(ns.Index(i), y.Index(i))
		}
		x.Set(ns)

	default:
		s.panic("setArray", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) setMap(x, y reflect.Value) {
	x = indirect(x)
	switch x.Kind() {
	case reflect.Interface:
		if x.NumMethod() != 0 {
			s.panic("setMap", setError{src: y.Type(), dst: x.Type()})
		}
		x.Set(y)

	case reflect.Map:
		nm := reflect.MakeMap(x.Type())
		for _, k := range y.MapKeys() {
			nm.SetMapIndex(k, y.MapIndex(k))
		}
		x.Set(nm)

	case reflect.Struct:
		if y.Type().Key().Kind() != reflect.String {
			s.panic("setMap", setError{src: y.Type(), dst: x.Type()})
		}
		fs := typeFields(x.Type())
		for _, k := range y.MapKeys() {
			f, ok := fs[k.String()]
			if !ok {
				continue
			}
			if !f.writeable {
				continue
			}
			s.set(x.Field(f.index), y.MapIndex(k))
		}

	default:
		s.panic("setMap", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) setStruct(x, y reflect.Value) {
	x = indirect(x)
	switch x.Kind() {
	case reflect.Interface:
		if x.NumMethod() != 0 {
			s.panic("setStruct", setError{src: y.Type(), dst: x.Type()})
		}
		x.Set(y)

	case reflect.Map:
		if x.Type().Key().Kind() != reflect.String {
			s.panic("setStruct", setError{src: y.Type(), dst: x.Type()})
		}
		nm := reflect.MakeMap(x.Type())
		fs := typeFields(y.Type())
		for _, f := range fs {
			nm.SetMapIndex(reflect.ValueOf(f.name), y.Field(f.index))
		}
		x.Set(nm)

	case reflect.Struct:
		xfs := typeFields(x.Type())
		yfs := typeFields(y.Type())
		for k, yf := range yfs {
			xf, ok := xfs[k]
			if !ok {
				continue
			}
			if !xf.writeable {
				continue
			}
			s.set(x.Field(xf.index), y.Field(yf.index))
		}

	default:
		s.panic("setStruct", setError{src: y.Type(), dst: x.Type()})
	}
}

func (s setState) panic(method string, err error) {
	panic(errs.ErrorAt(method, err))
}
