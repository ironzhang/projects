package resolver

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"

	"github.com/ironzhang/super-dns-go/superdns/pkg/model"
)

type provider struct {
	model atomic.Value // *model.ProviderModel
	route atomic.Value // RoutePolicy
}

func (p *provider) reload(m model.ProviderModel) {
	p.model.Store(&m)
}

func (p *provider) lookupCluster(ctx context.Context, tags map[string]string) (model.Cluster, error) {
	m, ok := p.model.Load().(*model.ProviderModel)
	if !ok {
		return model.Cluster{}, ErrInvalidProviderModel
	}
	r, ok := p.route.Load().(RoutePolicy)
	if !ok {
		return model.Cluster{}, ErrInvalidRoutePolicy
	}
	return lookupCluster(m, r, tags)
}

func lookupCluster(m *model.ProviderModel, r RoutePolicy, tags map[string]string) (model.Cluster, error) {
	dests := r.MatchRoute(tags, m.Clusters)
	if len(dests) <= 0 {
		dests = m.DefaultDestinations
	}

	cname, err := pick(dests)
	if err != nil {
		return model.Cluster{}, err
	}
	c, ok := m.Clusters[cname]
	if !ok {
		return model.Cluster{}, fmt.Errorf("%w: can not find %q cluster", ErrClusterNotFound, cname)
	}
	return c, nil
}

func pick(dests []model.Destination) (cluster string, err error) {
	sum := 0.0
	p := rand.Float64()
	for _, dest := range dests {
		sum += dest.Percent
		if p < sum {
			return dest.Cluster, nil
		}
	}
	if len(dests) > 0 {
		return dests[0].Cluster, nil
	}
	return "", ErrInvalidDestinations
}
