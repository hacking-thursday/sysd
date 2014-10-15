package mods

import (
	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/version"
	"net/http"
)

var (
	Modules = map[string]map[string]HttpApiFunc{
		"GET":     {},
		"POST":    {},
		"DELETE":  {},
		"OPTIONS": {},
	}
)

func init() {
	log.Debugf("pkg mods init()")
}

type HttpApiFunc func(eng interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) error

func Register(method string, route string, fct HttpApiFunc) error {
	if _, exists := Modules[method][route]; exists {
		log.Debugf("HttpApiFunc already registered %s::%s", method, route)
		return nil
	}

	Modules[method][route] = fct

	return nil
}
