package pb_parser

import (
	"os"
	"testing"
)

var metas = []Meta{
	{
		Name: "Point",
		Fields: []Field{
			{Name: "X", Type: "int32", Tag: 1},
			{Name: "Y", Type: "int32", Tag: 2},
		},
	},
	{
		Name: "Line",
		Fields: []Field{
			{Name: "Points", Type: "Point", Tag: 1, Repeat: true},
		},
	},
}

func TestMetaWriteHFile(t *testing.T) {
	for _, m := range metas {
		if err := m.WriteHFile(os.Stdout); err != nil {
			t.Fatal(err)
		}
	}
}

func TestMetaWriteCFile(t *testing.T) {
	for _, m := range metas {
		if err := m.WriteCFile(os.Stdout); err != nil {
			t.Fatal(err)
		}
	}
}
