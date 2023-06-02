package main

import (
	"github.com/ironzhang/super-dns-go/superdns/pkg/model"
	"github.com/ironzhang/super-dns-go/superdns/pkg/superutil"
)

func main() {
	m := model.ProviderModel{
		Domain: "www.superdns.com",
		Clusters: map[string]model.Cluster{
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
						Weight: 100,
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
						Weight: 100,
					},
				},
			},
		},
		DefaultDestinations: []model.Destination{
			{
				Cluster: "hna",
				Percent: 0.5,
			},
			{
				Cluster: "hnb",
				Percent: 0.5,
			},
		},
	}

	superutil.WriteJSON("www.superdns.com.json", m)
}
