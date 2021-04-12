package restful

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/httputils"
	"github.com/ironzhang/matrix/restful/codec"
	"github.com/ironzhang/matrix/tlog"
	"github.com/ironzhang/matrix/uuid"
)

func NewServeMux(c codec.Codec) *ServeMux {
	if c == nil {
		c = codec.DefaultCodec
	}
	return &ServeMux{
		codec:    c,
		patterns: make([]*pattern, 0),
	}
}

type ServeMux struct {
	codec    codec.Codec
	patterns []*pattern
}

func (m *ServeMux) Delete(pat string, h interface{}) error {
	return m.Add("DELETE", pat, h)
}

func (m *ServeMux) Get(pat string, h interface{}) error {
	return m.Add("GET", pat, h)
}

func (m *ServeMux) Head(pat string, h interface{}) error {
	return m.Add("HEAD", pat, h)
}

func (m *ServeMux) Options(pat string, h interface{}) error {
	return m.Add("OPTIONS", pat, h)
}

func (m *ServeMux) Patch(pat string, h interface{}) error {
	return m.Add("PATCH", pat, h)
}

func (m *ServeMux) Post(pat string, h interface{}) error {
	return m.Add("POST", pat, h)
}

func (m *ServeMux) Put(pat string, h interface{}) error {
	return m.Add("PUT", pat, h)
}

func (m *ServeMux) Add(meth, pat string, i interface{}) error {
	h, err := parseHandler(i)
	if err != nil {
		return fmt.Errorf("parse handler: %v", err)
	}

	for _, p := range m.patterns {
		if p.pat == pat {
			p.add(meth, h)
			return nil
		}
	}
	p := newPattern(pat)
	p.add(meth, h)
	m.patterns = append(m.patterns, p)
	return nil
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context_value.WithTraceId(context.Background(), getTraceId(r.Header))
	if err := m.serveHTTP(ctx, w, r); err != nil {
		m.setError(w, err)
	}
}

func (m *ServeMux) serveHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	log := tlog.WithContext(ctx).Sugar().With("method", r.Method, "path", r.URL.Path)

	var found bool
	for _, p := range m.patterns {
		v, ok := p.try(r.URL.Path)
		if !ok {
			continue
		}
		found = true
		h, ok := p.get(r.Method)
		if !ok {
			continue
		}
		return m.serve(ctx, h, v, w, r)
	}

	if found {
		log.Info(http.StatusText(http.StatusMethodNotAllowed))
		return Errorf(http.StatusMethodNotAllowed, codes.NotAllowed, "method(%s) not allowed", r.Method)
	} else {
		log.Info(http.StatusText(http.StatusNotFound))
		return Errorf(http.StatusNotFound, codes.NotFound, "page(%s) not found", r.URL.Path)
	}
}

func (m *ServeMux) serve(ctx context.Context, h *handler, v url.Values, w http.ResponseWriter, r *http.Request) (err error) {
	log := tlog.WithContext(ctx).Sugar().With("method", r.Method, "path", r.URL.Path)

	// check Content-Type
	if err = m.checkContentType(r.Header); err != nil {
		log.Infow("check content type", "error", err)
		return Errorf(http.StatusBadRequest, codes.InvalidHeader, err.Error())
	}

	args := newReflectValue(h.args)
	reply := newReflectValue(h.reply)

	// Decode
	if !isNilInterface(h.args) {
		if err = m.codec.Decode(r.Body, args.Interface()); err != nil {
			log.Infow("decode", "error", err)
			return Errorf(http.StatusBadRequest, codes.DecodeFail, err.Error())
		}
	}

	// with context
	ctx = context_value.WithRequest(ctx, r)
	ctx = context_value.WithResponseWriter(ctx, w)

	// Handle
	if err = h.Handle(ctx, v, args, reply); err != nil {
		log.Infow("handle", "error", err)
		return err
	}

	// Encode
	if !isNilInterface(h.reply) {
		w.Header().Set("Content-Type", m.codec.ContentType())
		if err = m.codec.Encode(w, reply.Interface()); err != nil {
			log.Errorw("encode", "error", err)
			return Errorf(http.StatusInternalServerError, codes.EncodeFail, err.Error())
		}
	}

	return nil
}

func (m *ServeMux) checkContentType(h http.Header) error {
	if v := h.Get("Content-Type"); v != "" && v != m.codec.ContentType() {
		return fmt.Errorf("Content-Type not %s: %s", m.codec.ContentType(), v)
	}
	return nil
}

func (m *ServeMux) setError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if te, ok := err.(HTTPStatus); ok {
		status = te.HTTPStatus()
	}
	e := codec.ToError(err)
	w.Header().Set("Content-Type", m.codec.ContentType())
	w.WriteHeader(status)
	m.codec.EncodeError(w, e)
}

func getTraceId(h http.Header) string {
	if v := h.Get(httputils.X_TRACE_ID); v != "" {
		return v
	}
	return uuid.New().String()
}
