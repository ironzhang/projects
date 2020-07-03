package resolver

import (
	"errors"
	"strings"

	gr "google.golang.org/grpc/resolver"
)

type resolver struct {
	cc gr.ClientConn
}

func (p *resolver) start(target gr.Target) error {
	if target.Authority != "seeds" {
		return errors.New("invalid target")
	}
	endpoints := strings.Split(target.Endpoint, ",")
	if len(endpoints) <= 0 {
		return errors.New("invalid target")
	}

	addrs := make([]gr.Address, 0, len(endpoints))
	for _, endpoint := range endpoints {
		addrs = append(addrs, gr.Address{Addr: endpoint})
	}
	p.cc.UpdateState(gr.State{Addresses: addrs})

	return nil
}

func (p *resolver) ResolveNow(gr.ResolveNowOptions) {
}

func (p *resolver) Close() {
}

type builder struct{}

func (builder) Build(target gr.Target, cc gr.ClientConn,
	opts gr.BuildOptions) (gr.Resolver, error) {
	r := &resolver{
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
	gr.Register(builder{})
}
