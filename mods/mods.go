package mods

import (
	"fmt"
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

type HttpApiFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error)

func Register(method string, route string, fct HttpApiFunc) (err error) {
	if _, exists := Modules[method][route]; exists {
		err = fmt.Errorf("HttpApiFunc already registered %s::%s", method, route)
		return
	}

	Modules[method][route] = fct
	return
}
