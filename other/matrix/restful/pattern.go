package restful

import (
	"net/url"
	"strings"
)

type pattern struct {
	pat      string
	handlers map[string]*handler
}

func newPattern(pat string) *pattern {
	return &pattern{pat: pat, handlers: make(map[string]*handler)}
}

func (p *pattern) add(meth string, h *handler) {
	p.handlers[strings.ToUpper(meth)] = h
}

func (p *pattern) get(meth string) (*handler, bool) {
	h, ok := p.handlers[strings.ToUpper(meth)]
	return h, ok
}

func (p *pattern) try(path string) (url.Values, bool) {
	v := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(p.pat):
			if p.pat != "/" && len(p.pat) > 0 && p.pat[len(p.pat)-1] == '/' {
				return v, true
			}
			return nil, false
		case p.pat[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(p.pat, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			v.Add(":"+name, val)
		case path[i] == p.pat[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(p.pat) {
		return nil, false
	}
	return v, true
}

func matchPart(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
