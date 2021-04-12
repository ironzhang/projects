package httputils

import (
	"context"
	"net"
	"net/http"
)

type HTTPServer struct {
	l net.Listener
	h http.Handler
}

func (s *HTTPServer) Init(addr string, h http.Handler) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.l = l
	s.h = h
	return nil
}

func (s *HTTPServer) Serve(ctx context.Context) {
	Serve(ctx, s.l, s.h)
}

func Serve(ctx context.Context, l net.Listener, h http.Handler) {
	s := http.Server{Handler: h}
	go func() {
		<-ctx.Done()
		s.Close()
	}()
	s.Serve(l)
}
