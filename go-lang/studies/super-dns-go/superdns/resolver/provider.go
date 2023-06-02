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
}

func (p *provider) reload(m model.ProviderModel) {
	p.model.Store(&m)
}

func (p *provider) lookupCluster(ctx context.Context, tags map[string]string) (model.Cluster, error) {
	m, ok := p.model.Load().(*model.ProviderModel)
	if !ok {
		return model.Cluster{}, ErrInvalidProviderModel
	}
	return lookupCluster(m, tags)
}

func lookupCluster(m *model.ProviderModel, tags map[string]string) (model.Cluster, error) {
	cname, err := pick(m.DefaultDestinations)
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
