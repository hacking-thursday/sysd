package mods

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
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

type HttpApiFunc func(eng interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) error

func Register(method string, route string, fct HttpApiFunc) (err error) {
	if _, exists := Modules[method][route]; exists {
		err = fmt.Errorf("HttpApiFunc already registered %s::%s", method, route)
		return
	}

	Modules[method][route] = fct
	return
}
