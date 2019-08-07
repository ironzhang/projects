package api

import "github.com/ironzhang/matrix/rest"

func Register(r *rest.Rest) error {
	r.Post("account")
	r.Delete("account/:uid")
}
