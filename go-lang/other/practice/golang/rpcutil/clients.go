package rpcutil

import (
	"net/rpc"
	"runtime"
	"time"

	"github.com/ironzhang/golang/lru"
)

type ClientDialer func(network, address string) (*rpc.Client, error)

type Clients struct {
	network string
	dial    ClientDialer
	cache   *lru.Cache
}

func NewClients(network string, dial ClientDialer, maxcache int) *Clients {
	if dial == nil {
		dial = rpc.Dial
	}
	return &Clients{
		network: network,
		dial:    dial,
		cache:   lru.New(maxcache, nil),
	}
}

func (c *Clients) Call(addr string, method string, args, reply interface{}, timeout time.Duration) error {
	if v, ok := c.cache.Get(addr); ok {
		cli := v.(*rpc.Client)
		return Call(cli, method, args, reply, timeout)
	}
	cli, err := c.dial(c.network, addr)
	if err != nil {
		return err
	}
	runtime.SetFinalizer(cli, (*rpc.Client).Close)
	c.cache.Add(addr, cli)
	return Call(cli, method, args, reply, timeout)
}
