package luavm

import (
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

func GetFunctionFromTable(t *lua.LTable, key string) (*lua.LFunction, error) {
	value := t.RawGetString(key)
	fn, ok := value.(*lua.LFunction)
	if !ok {
		return nil, fmt.Errorf("table[%s] is not a function", key)
	}
	return fn, nil
}

func ParseValue(a interface{}, lv lua.LValue) error {
	return nil
}

func parseValue(lv lua.LValue, rv reflect.Value) {
	switch k := rv.Type().Kind(); k {
	case reflect.Ptr:
	case reflect.Interface:
	case reflect.Array:
	case reflect.Slice:
	case reflect.Map:
	case reflect.Struct:
	default:
	}
}

func parsePtrValue(lv lua.LValue, rv reflect.Value) {
	e := reflect.New(rv.Type().Elem())
	parseValue(lv, e)
	rv.Elem().Set(e)
}

func parseInterfaceValue(lv lua.LValue, rv reflect.Value) {
	e := rv.Elem()
	if !e.IsValid() {
		return
	}
	parseValue(lv, e)
}

func parseArrayValue(lv lua.LValue, rv reflect.Value) {
	lt := lv.(*lua.LTable)
	for i := 0; i < lt.Len(); i++ {
		parseValue(lt.RawGetInt(i), rv.Index(i))
	}
}

func parseMapValue(lv lua.LValue, rv reflect.Value) {
	rv.Elem().Set(reflect.MakeMap(rv.))

	lt.ForEach(func(key lua.LValue, value lua.LValue) {
		rv.Type()

		parseValue(key)
		parseValue(value)

		rv.SetMapIndex()
	})
}

func parseBaseValue(lv lua.LValue, k reflect.Kind, rv reflect.Value) {
	switch k {
	case reflect.Bool:
		rv.SetBool(lua.LVAsBool(lv))

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		rv.SetInt(int64(lua.LVAsNumber(lv)))

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		rv.SetUint(uint64(lua.LVAsNumber(lv)))

	case reflect.Float32,
		reflect.Float64:
		rv.SetFloat(float64(lua.LVAsNumber(lv)))

	case reflect.String:
		rv.SetString(lua.LVAsString(lv))

	default:
		panic(fmt.Errorf("%s is an unsupported base type", k.String()))
	}
}
