package mods

import (
	"io"
	"net/http"

	"github.com/docker/docker/pkg/log"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/gorilla/mux"
	"github.com/tsaikd/KDGoLib/env"
)

var (
	flApiPrefix = flag.String(
		[]string{"-SYSD_API_PREFIX"},
		env.GetString("SYSD_API_PREFIX", ""),
		"Sysd API Server URL Prefix",
	)
)

func CreateRouter() (r *mux.Router, err error) {
	var (
		prefix = *flApiPrefix
	)
	r = mux.NewRouter()

	for method, routes := range Modules {
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

// used for testing
func NewApiRequest(method string, urlStr string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, *flApiPrefix+urlStr, body)
	return
}
