package balance

import (
	"errors"

	"github.com/ironzhang/sidecar/model"
)

func SWRR(addrs []*model.Address) (*model.Address, error) {
	var best *model.Address

	total := 0
	for _, node := range addrs {
		if best == nil || node.Current > best.Current {
			best = node
		}
		total += node.Current
		node.Current += node.Weight
	}
	if best == nil {
		return nil, errors.New("no node")
	}
	best.Current -= total

	return best, nil
}
