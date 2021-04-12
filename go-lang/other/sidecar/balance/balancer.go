package balance

import (
	"context"

	"github.com/ironzhang/sidecar/model"
)

type Balancer interface {
	SelectAddress(ctx context.Context, addrs []*model.Address) (*model.Address, error)
}
