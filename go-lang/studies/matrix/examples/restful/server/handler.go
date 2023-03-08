package main

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ironzhang/matrix/restful"
)

func NewHTTPHandler() (http.Handler, error) {
	var err error
	var h handlers
	m := restful.NewServeMux(nil)
	//m.SetVerbose(2)
	if err = m.Get("/", h.Root); err != nil {
		return nil, err
	}
	if err = m.Post("/echo", h.Echo); err != nil {
		return nil, err
	}
	if err = m.Get("/hello/:name", h.Hello); err != nil {
		return nil, err
	}
	return m, nil
}

type handlers struct{}

func (h *handlers) Root(ctx context.Context, values url.Values, req interface{}, resp *string) error {
	*resp = "hello, restful"
	return nil
}

func (h *handlers) Echo(ctx context.Context, values url.Values, req string, resp *string) error {
	*resp = req
	return nil
}

func (h *handlers) Hello(ctx context.Context, values url.Values, req interface{}, resp *string) error {
	*resp = values.Get(":name")
	return nil
}
