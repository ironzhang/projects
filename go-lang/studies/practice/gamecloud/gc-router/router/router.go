package router

import "net/http"

type Router struct {
}

func (p *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
}
