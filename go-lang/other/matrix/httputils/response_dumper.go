package httputils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func NewResponseDumper(w http.ResponseWriter, r *http.Request) *ResponseDumper {
	return &ResponseDumper{
		ResponseWriter: w,
		proto:          r.Proto,
		status:         http.StatusOK,
	}
}

type ResponseDumper struct {
	http.ResponseWriter
	proto  string
	status int
	buffer bytes.Buffer
}

func (p *ResponseDumper) WriteHeader(status int) {
	p.status = status
	p.ResponseWriter.WriteHeader(status)
}

func (p *ResponseDumper) Write(b []byte) (int, error) {
	p.buffer.Write(b)
	return p.ResponseWriter.Write(b)
}

func (p *ResponseDumper) Dump(body bool) []byte {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%s %d %s\r\n", p.proto, p.status, http.StatusText(p.status))
	if len(p.Header()) > 0 {
		p.Header().Write(&out)
	}
	fmt.Fprintf(&out, "\r\n")
	if body && p.buffer.Len() > 0 {
		io.Copy(&out, &p.buffer)
	}
	return out.Bytes()
}
