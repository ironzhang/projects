package backend

import (
	"context"

	"github.com/ironzhang/sidecar/model"
)

type Pool struct {
}

func (p *Pool) SelectAddress(ctx context.Context, tags model.Tags) (model.Address, error) {
	return "", nil
}

func (p *Pool) MarkAddressState(ctx context.Context, addr model.Address, fault bool) {
}
