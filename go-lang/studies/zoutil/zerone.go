package zoutil

import (
	"sync"
	"time"

	"github.com/ironzhang/zerone"
	"github.com/ironzhang/zerone/zclient"
	"github.com/ironzhang/zerone/zserver"
)

type Zerone struct {
	Timeout time.Duration

	node          string
	clientVerbose int
	serverVerbose int

	zerone  zerone.Zerone
	clients sync.Map
	servers sync.Map
}

func Open(node string, opts Options) (*Zerone, error) {
	zopts, err := opts.zerone()
	if err != nil {
		return nil, err
	}
	z, err := zerone.NewZerone(zopts)
	if err != nil {
		return nil, err
	}
	return &Zerone{
		Timeout:       time.Duration(opts.Timeout),
		node:          node,
		clientVerbose: opts.ClientVerbose,
		serverVerbose: opts.ServerVerbose,
		zerone:        z,
	}, nil
}

func (z *Zerone) Close() error {
	z.clients.Range(func(key, value interface{}) bool {
		c := value.(*zclient.Client)
		c.Close()
		return true
	})

	z.servers.Range(func(key, value interface{}) bool {
		s := value.(*zserver.Server)
		s.Close()
		return true
	})

	return z.zerone.Close()
}

func (z *Zerone) Node() string {
	return z.node
}

func (z *Zerone) GetClientVerbose() int {
	return z.clientVerbose
}

func (z *Zerone) SetClientVerbose(verbose int) {
	z.clientVerbose = verbose
	z.clients.Range(func(key, value interface{}) bool {
		c := value.(*zclient.Client)
		c.SetTraceVerbose(verbose)
		return true
	})
}

func (z *Zerone) GetServerVerbose() int {
	return z.serverVerbose
}

func (z *Zerone) SetServerVerbose(verbose int) {
	z.serverVerbose = verbose
	z.servers.Range(func(key, value interface{}) bool {
		s := value.(*zserver.Server)
		s.SetTraceVerbose(verbose)
		return true
	})
}

func (z *Zerone) NewClient(service string) (*zclient.Client, error) {
	if v, ok := z.clients.Load(service); ok {
		return v.(*zclient.Client), nil
	}

	c, err := z.zerone.NewClient(z.node, service)
	if err != nil {
		return nil, err
	}
	c.SetTraceVerbose(z.clientVerbose)

	if v, loaded := z.clients.LoadOrStore(service, c); loaded {
		c.Close()
		return v.(*zclient.Client), nil
	}
	return c, nil
}

func (z *Zerone) GetClient(service string) (*zclient.Client, bool) {
	if v, ok := z.clients.Load(service); ok {
		return v.(*zclient.Client), true
	}
	return nil, false
}

func (z *Zerone) NewServer(service string) (*zserver.Server, error) {
	if v, ok := z.servers.Load(service); ok {
		return v.(*zserver.Server), nil
	}

	s, err := z.zerone.NewServer(z.node, service)
	if err != nil {
		return nil, err
	}
	s.SetTraceVerbose(z.serverVerbose)

	if v, loaded := z.servers.LoadOrStore(service, s); loaded {
		s.Close()
		return v.(*zserver.Server), nil
	}
	return s, nil
}

func (z *Zerone) GetServer(service string) (*zserver.Server, bool) {
	if v, ok := z.servers.Load(service); ok {
		return v.(*zserver.Server), true
	}
	return nil, false
}
