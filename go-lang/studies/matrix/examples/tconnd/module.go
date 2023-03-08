package main

import (
	"context"
	"net"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/netutils/listen-mux"
	"github.com/ironzhang/matrix/tlog"
)

var Config = &C{
	Addrs: []string{":1883", ":2000", ":2001", ":2002", ":2003", ":2004", ":2005", ":2006", ":2007", ":2008"},
}

var Module = &M{}

func init() {
	framework.Register(Module, nil, Config)
}

type C struct {
	Addrs []string `json:",readonly"`
}

type M struct {
	ln net.Listener
}

func (m *M) Name() string {
	return "main-module"
}

func (m *M) Init() (err error) {
	m.ln, err = listen_mux.Listen("tcp", Config.Addrs, 0)
	if err != nil {
		return err
	}
	return nil
}

func (m *M) Fini() error {
	return nil
}

func (m *M) Run(ctx context.Context) {
	log := tlog.Std().Sugar().With("module", m.Name(), "addr", m.ln.Addr().String())
	log.Info("start")
	serve(ctx, m.ln)
	log.Info("stop")
}
