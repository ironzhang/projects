package main

import (
	"context"
	"time"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/super-dns-go/superdns/resolver"
)

func main() {
	r := resolver.Resolver{
		ResourcePath:  "./superdns",
		WatchInterval: time.Second,
	}

	for {
		c, err := r.LookupCluster(context.Background(), "www.superdns.com", nil)
		if err != nil {
			tlog.Errorw("lookup cluster", "error", err)
		} else {
			tlog.Infow("lookup cluster", "cluster", c)
		}
		time.Sleep(10 * time.Second)
	}
}
