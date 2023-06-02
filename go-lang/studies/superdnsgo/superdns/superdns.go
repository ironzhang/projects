package superdns

import (
	"context"
)

type Endpoint struct {
	IP     string
	Port   int
	Weight int
	Status string
}

type Cluster struct {
	Name      string
	Endpoints []Endpoint
}

func SetLoadBalancer() {
}

func LookupEndpoint(ctx context.Context, domain string, tags map[string]string) (Endpoint, error) {
}

func LookupCluster(ctx context.Context, domain string, tags map[string]string) (Cluster, error) {
}
