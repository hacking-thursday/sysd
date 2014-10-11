package server

import (
	"github.com/docker/docker/pkg/log"
	"github.com/gorilla/mux"

	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	//mods.Register("GET", "/ping", ping)
	//mods.Register("GET", "/memstats", memstats)
        // will be loaded in /mods/ subdirectories
}

func createRouter() (r *mux.Router, err error) {
	var (
		prefix = *flApiPrefix
	)
	r = mux.NewRouter()

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
