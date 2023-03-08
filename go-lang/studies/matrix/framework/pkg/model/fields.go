package model

import (
	"reflect"

	"github.com/ironzhang/matrix/framework/pkg/tags"
)

func parseTag(tag string) (string, bool) {
	if tag != "" {
		name, opts := tags.ParseTag(tag)
		if tags.IsValidTag(name) {
			return name, opts.Contains("writeable")
		}
		return "", opts.Contains("writeable")
	}
	return "", false
}

type field struct {
	name      string
	index     int
	typ       reflect.Type
	writeable bool
}

func typeFields(t reflect.Type) map[string]field {
	var name string
	var writeable bool
	fields := make(map[string]field)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}
		name, writeable = parseTag(sf.Tag.Get("json"))
		if name == "-" {
			continue
		}
		if name == "" {
			name = sf.Name
		}
		fields[name] = field{name: name, index: i, typ: sf.Type, writeable: writeable}
	}
	return fields
}
