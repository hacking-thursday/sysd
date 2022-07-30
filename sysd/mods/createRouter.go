package mods

import (
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "flag"
	"github.com/gorilla/mux"
	"github.com/tsaikd/KDGoLib/env"
)

var (
	flApiPrefix = flag.String(
		"SYSD_API_PREFIX",
		env.GetString("SYSD_API_PREFIX", ""),
		"Sysd API Server URL Prefix",
	)
)

func CreateRouter(eng interface{}) (r *mux.Router, err error) {
	var (
		prefix = strings.TrimSuffix(*flApiPrefix, "/")
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
			f := makeHttpHandler(eng, localMethod, localRoute, localFct)

			if strings.HasSuffix(localRoute, "/*") {
				routeBase := prefix + "/" + strings.TrimSuffix(localRoute, "/*")
				routeBase = strings.Replace(routeBase, "//", "/", -1)
				r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
					return strings.HasPrefix(r.URL.Path, routeBase)
				}).HandlerFunc(f)
			} else if prefix == "" {
				r.Path(localRoute).Methods(localMethod).HandlerFunc(f)
			} else {
				r.PathPrefix(prefix).Path(localRoute).Methods(localMethod).HandlerFunc(f)
			}
		}
	}

	return
}

func makeHttpHandler(eng interface{}, localMethod string, localRoute string, handlerFunc HttpApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log the request
		log.Debugf("Calling %s %s", localMethod, localRoute)
		log.Infof("%s %s", r.Method, r.RequestURI)

		writeCorsHeaders(w, r)

		if err := handlerFunc(eng, w, r, mux.Vars(r)); err != nil {
			log.Errorf("Handler for %s %s returned error: %s", localMethod, localRoute, err)
			HttpError(w, err)
		}
	}
}

func writeCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

// used for testing
func NewApiRequest(method string, urlStr string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, *flApiPrefix+urlStr, body)
	return
}
