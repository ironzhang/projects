package resolver

import (
	"context"
	"errors"
	"testing"

	"github.com/ironzhang/super-dns-go/superdns/pkg/model"
)

var TestClusters = map[string]model.Cluster{
	"hna": model.Cluster{
		Name: "hna",
		Endpoints: []model.Endpoint{
			{
				Addr:   "127.0.0.1:8000",
				State:  model.Enabled,
				Weight: 100,
			},
			{
				Addr:   "127.0.0.2:8000",
				State:  model.Enabled,
				Weight: 200,
			},
		},
	},
	"hnb": model.Cluster{
		Name: "hnb",
		Endpoints: []model.Endpoint{
			{
				Addr:   "128.0.0.1:8000",
				State:  model.Enabled,
				Weight: 100,
			},
			{
				Addr:   "128.0.0.2:8000",
				State:  model.Enabled,
				Weight: 200,
			},
		},
	},
}

func TestProviderLookupCluster(t *testing.T) {
	tests := []struct {
		model  model.ProviderModel
		target string
		err    error
	}{
		{
			model: model.ProviderModel{
				Domain:   "www.abc.com",
				Clusters: TestClusters,
				DefaultDestinations: []model.Destination{
					{Cluster: "hna", Percent: 1},
				},
			},
			target: "hna",
		},
		{
			model: model.ProviderModel{
				Domain:   "www.abc.com",
				Clusters: TestClusters,
				DefaultDestinations: []model.Destination{
					{Cluster: "hnb", Percent: 1},
				},
			},
			target: "hnb",
		},
		{
			model: model.ProviderModel{
				Domain:   "www.abc.com",
				Clusters: TestClusters,
				DefaultDestinations: []model.Destination{
					{Cluster: "hnc", Percent: 1},
				},
			},
			err: ErrClusterNotFound,
		},
	}
	for i, tt := range tests {
		var p provider
		p.reload(tt.model)
		cluster, err := p.lookupCluster(context.Background(), nil)
		if !errors.Is(err, tt.err) {
			t.Fatalf("%d: error match: got '%v', want '%v'", i, err, tt.err)
		}
		if err != nil {
			t.Logf("%d: lookup cluster: %v", i, err)
			continue
		}
		if got, want := cluster.Name, tt.target; got != want {
			t.Fatalf("%d: cluster match: got %q, want %q", i, got, want)
		}
		t.Logf("cluster: %v", cluster)
	}
}
