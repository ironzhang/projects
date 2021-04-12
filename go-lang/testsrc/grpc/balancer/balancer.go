package balancer

import (
	"log"
	"math/rand"
	"sync"

	gbl "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Name = "conn_round_robin"

type picker struct {
	sc gbl.SubConn
}

func (p *picker) Pick(gbl.PickInfo) (gbl.PickResult, error) {
	return gbl.PickResult{SubConn: p.sc}, nil
}

type pickerBuilder struct {
	mu   sync.Mutex
	last gbl.SubConn
}

func (p *pickerBuilder) Build(info base.PickerBuildInfo) gbl.V2Picker {
	log.Printf("pickerBuilder: new picker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(gbl.ErrNoSubConnAvailable)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.last != nil {
		if _, ok := info.ReadySCs[p.last]; ok {
			return &picker{sc: p.last}
		}
	}

	scs := make([]gbl.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	p.last = scs[rand.Intn(len(scs))]
	return &picker{sc: p.last}
}

func newBalancerBuilder() gbl.Builder {
	return base.NewBalancerBuilderV2(Name, &pickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	gbl.Register(newBalancerBuilder())
}
