package server

import (
	"github.com/docker/docker/pkg/log"
	"github.com/gorilla/mux"

	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	mods.Register("GET", "/osver", osver)
}

func createRouter_windows(r *mux.Router) {
	var (
		prefix = *flApiPrefix
	)

	for method, routes := range mods.Modules {
		for route, fct := range routes {
			log.Debugf("Registering %s, %s", method, route)
			// NOTE: scope issue, make sure the variables are local and won't be changed
			localRoute := route
			localFct := fct
			localMethod := method

			// build the handler function
			f := makeHttpHandler(localMethod, localRoute, localFct)

			if prefix == "" {
				r.Path(localRoute).Methods(localMethod).HandlerFunc(f)
			} else {
				r.PathPrefix(prefix).Path(localRoute).Methods(localMethod).HandlerFunc(f)
			}
		}
	}

	return
}
