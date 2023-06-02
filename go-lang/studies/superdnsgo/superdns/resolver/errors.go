package resolver

import "errors"

var (
	ErrInvalidProviderModel = errors.New("invalid provider model")
	ErrInvalidRoutePolicy   = errors.New("invalid route policy")
	ErrInvalidDestinations  = errors.New("invalid destinations")
	ErrClusterNotFound      = errors.New("can not find cluster")
)
