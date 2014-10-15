package apilist

import (
	"net/http"

	"github.com/docker/docker/pkg/version"

	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	mods.Register("GET", "/apilist", apilist)
}

func apilist(engine interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var (
		out       []byte
		moduleMap = map[string]bool{}
		modules   = []string{}
	)

	for _, routes := range mods.Modules {
		for route, _ := range routes {
			moduleMap[route] = true
		}
	}

	for route, _ := range moduleMap {
		modules = append(modules, route)
	}

	if out, err = mods.Marshal(r, modules); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
