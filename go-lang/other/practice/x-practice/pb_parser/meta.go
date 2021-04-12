package pb_parser

import (
	"io"
	"text/template"
)

type Field struct {
	Name   string
	Type   string
	Tag    int
	Repeat bool
}

type Meta struct {
	Name   string
	Fields []Field
}

func (m *Meta) WriteHFile(w io.Writer) error {
	tmpl, err := template.New("").Parse(hfileTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, m)
}

func (m *Meta) WriteCFile(w io.Writer) error {
	tmpl, err := template.New("").Parse(cfileTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, m)
}

var hfileTemplate = `
typedef struct {{.Name}} {
	{{range .Fields}} {{if .Repeat}} vertor<{{.Type}}> {{.Name}}; {{else}} {{.Type}} {{.Name}}; {{end}}
	{{end}}
}{{.Name}};
`

var cfileTemplate = `
inline size_t Encode{{.Name}}({{.Name}} *p, uint8 *buf, size_t sz, int tag) {
	size_t n = 0;
	size_t idx = 0;

	n = EncodeVarint(tag, buf, sz);
	if (n < 0)
		return n;
	idx += n;
	sz -= n;
	{{range .Fields}} {{if .Repeat}}
	int64 length = int64(p->{{.Name}}.size());
	n = EncodeVarint(&length, buf+idx, sz);
	if (n < 0)
		return n;
	idx += n;
	sz -= n;

	for (int64 i := 0; i < length; i++) {
		n = Encode{{.Type}}(&p->{{.Name}}[i], buf+idx, sz, {{.Tag}});
		if (n < 0)
			return n;
		idx += n;
		sz -= n;
	} {{else}}
	n = Encode{{.Type}}(&p->{{.Name}}, buf+idx, sz, {{.Tag}});
	if (n < 0)
		return n;
	idx += n;
	sz -= n; {{end}}
	{{end}}
	return s;
}
`
