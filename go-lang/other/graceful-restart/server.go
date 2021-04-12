package main

import (
	"net"
	"net/http"
	"os"

	"github.com/ironzhang/projects/other/graceful-restart/graceful"
	"github.com/ironzhang/tlog"
)

type Server struct {
	listener *net.TCPListener
}

func NewServer(addr string, fd uintptr) (*Server, error) {
	ln, err := graceful.ListenTCP("tcp", addr, fd)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: ln,
	}, nil
}

func (s *Server) Serve() {
	hs := http.Server{}
	hs.Serve(s.listener)
	hs.Close()
	hs.Shutdown()
}

func (s *Server) Shutdown() error {
	err := graceful.StopTCPListeners(s.listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Restart() error {
	pid, err := graceful.ForkExec(s.listener)
	if err != nil {
		tlog.Errorw("fork exec", "error", err)
		return err
	}
	tlog.Infow("fork exec success", "pid", pid)
	os.Exit(0)
	return nil
}
