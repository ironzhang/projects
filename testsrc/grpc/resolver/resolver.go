package resolver

import (
	"errors"
	"strings"

	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	cc resolver.ClientConn
}

func (p *Resolver) start(target resolver.Target) error {
	endpoints := strings.Split(target.Endpoint, ",")
	if len(endpoints) <= 0 {
		return errors.New("invalid target")
	}

	addrs := make([]resolver.Address, 0, len(endpoints))
	for _, endpoint := range endpoints {
		addrs = append(addrs, resolver.Address{Addr: endpoint})
	}
	p.cc.UpdateState(resolver.State{Addresses: addrs})

	return nil
}

func (p *Resolver) ResolveNow(resolver.ResolveNowOptions) {
}

func (p *Resolver) Close() {
}

type builder struct{}

func (builder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		cc: cc,
	}
	if err := r.start(target); err != nil {
		return nil, err
	}
	return r, nil
}

func (builder) Scheme() string {
	return "registry"
}

func init() {
	resolver.Register(builder{})
}
