package luavm

import (
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

type LVM struct {
	lstate *lua.LState
}

func New() *LVM {
	return &LVM{lstate: lua.NewState()}
}

func (p *LVM) Load(path string) error {
	return p.lstate.DoFile(path)
}

func (p *LVM) Call(fn *lua.LFunction, nret int, args ...lua.LValue) ([]lua.LValue, error) {
	err := p.lstate.CallByParam(lua.P{
		Fn:      fn,
		NRet:    nret,
		Protect: true,
	}, args...)
	if err != nil {
		return nil, err
	}

	var rets []lua.LValue
	for i := 1; i <= nret; i++ {
		rets = append(rets, p.lstate.Get(-1*i))
	}
	p.lstate.Pop(nret)

	return rets, nil
}

func (p *LVM) GetGlobal(name string) lua.LValue {
	return p.lstate.GetGlobal(name)
}

func (p *LVM) GetGlobalTable(name string) (*lua.LTable, error) {
	value := p.lstate.GetGlobal(name)
	t, ok := value.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("%q is not a table", name)
	}
	return t, nil
}

func (p *LVM) NewValue(a interface{}) (lv lua.LValue, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	v := reflect.ValueOf(a)
	return p.newValue(v), nil
}

func (p *LVM) newValue(v reflect.Value) lua.LValue {
	switch k := v.Type().Kind(); k {
	case reflect.Ptr, reflect.Interface:
		return p.newPtrOrInterfaceValue(v)
	case reflect.Array, reflect.Slice:
		return p.newArrayOrSliceValue(v)
	case reflect.Map:
		return p.newMapValue(v)
	case reflect.Struct:
		return p.newStructValue(v)
	default:
		return newBaseValue(k, v)
	}
}

func (p *LVM) newPtrOrInterfaceValue(v reflect.Value) lua.LValue {
	e := v.Elem()
	if !e.IsValid() {
		return lua.LNil
	}
	return p.newValue(e)
}

func (p *LVM) newArrayOrSliceValue(v reflect.Value) lua.LValue {
	lt := p.lstate.NewTable()
	for i := 0; i < v.Len(); i++ {
		lv := p.newValue(v.Index(i))
		lt.Append(lv)
	}
	return lt
}

func (p *LVM) newMapValue(v reflect.Value) lua.LValue {
	lt := p.lstate.NewTable()
	it := v.MapRange()
	for it.Next() {
		lk := p.newValue(it.Key())
		lv := p.newValue(it.Value())
		lt.RawSet(lk, lv)
	}
	return lt
}

func (p *LVM) newStructValue(v reflect.Value) lua.LValue {
	lt := p.lstate.NewTable()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		k := t.Field(i).Name
		lv := p.newValue(v.Field(i))
		lt.RawSetString(k, lv)
	}
	return lt
}

func newBaseValue(k reflect.Kind, v reflect.Value) lua.LValue {
	switch k {
	case reflect.Bool:
		return lua.LBool(v.Bool())

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return lua.LNumber(v.Int())

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return lua.LNumber(v.Uint())

	case reflect.Float32,
		reflect.Float64:
		return lua.LNumber(v.Float())

	case reflect.String:
		return lua.LString(v.String())

	default:
		panic(fmt.Errorf("%s is an unsupported base type", k.String()))
	}
}
