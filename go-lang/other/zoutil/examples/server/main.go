package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/ironzhang/x-pearls/log"
	"github.com/ironzhang/zerone/examples/rpc/arith"
	"github.com/ironzhang/zoutil"
	"github.com/ironzhang/zoutil/examples/conf"
)

type Options struct {
	net  string
	addr string
}

func (o *Options) Parse() {
	flag.StringVar(&o.net, "net", "tcp", "network")
	flag.StringVar(&o.addr, "addr", ":10000", "address")
	flag.Parse()
}

func main() {
	var opts Options
	opts.Parse()

	err := conf.LoadConfig("../conf/conf.toml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	z, err := zoutil.Open(conf.G.Node, conf.G.Zerone)
	if err != nil {
		log.Fatalf("zoutil open: %v", err)
	}
	defer z.Close()

	svr, err := z.NewServer("arith")
	if err != nil {
		log.Fatalf("new server: %v", err)
	}
	if err = svr.Register(arith.Arith(0)); err != nil {
		log.Fatalf("register: %v", err)
	}

	go func() {
		log.Infof("listen and serve on %s://%s", opts.net, opts.addr)
		if err = svr.ListenAndServe(opts.net, opts.addr, ""); err != nil {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
