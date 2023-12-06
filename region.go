package infra

import (
	"net/http"
)

type Region struct {
	Deployments map[string]*Deployment
}

func (region *Region) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := region.Deployments[r.Host]
	d.ServeHTTP(w, r)
}
