package resolver

import "errors"

var (
	ErrInvalidProviderModel = errors.New("invalid provider model")
	ErrInvalidDestinations  = errors.New("invalid destinations")
	ErrClusterNotFound      = errors.New("can not find cluster")
)
