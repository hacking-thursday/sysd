package mods

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/docker/docker/pkg/version"
	"github.com/gorilla/mux"
	"github.com/tsaikd/KDGoLib/env"
)

const (
	APIVERSION version.Version = "0.1"
)

var (
	flApiPrefix = flag.String(
		[]string{"-SYSD_API_PREFIX"},
		env.GetString("SYSD_API_PREFIX", ""),
		"Sysd API Server URL Prefix",
	)
)

func CreateRouter(eng interface{}) (r *mux.Router, err error) {
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
			f := makeHttpHandler(eng, localMethod, localRoute, localFct, APIVERSION)

			if prefix == "" {
				r.Path(localRoute).Methods(localMethod).HandlerFunc(f)
			} else {
				r.PathPrefix(prefix).Path(localRoute).Methods(localMethod).HandlerFunc(f)
			}
		}
	}

	return
}

func makeHttpHandler(eng interface{}, localMethod string, localRoute string, handlerFunc HttpApiFunc, dockerVersion version.Version) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log the request
		log.Debugf("Calling %s %s", localMethod, localRoute)
		log.Infof("%s %s", r.Method, r.RequestURI)

		if strings.Contains(r.Header.Get("User-Agent"), "Docker-Client/") {
			userAgent := strings.Split(r.Header.Get("User-Agent"), "/")
			if len(userAgent) == 2 && !dockerVersion.Equal(version.Version(userAgent[1])) {
				log.Debugf("Warning: client and server don't have the same version (client: %s, server: %s)", userAgent[1], dockerVersion)
			}
		}
		version := version.Version(mux.Vars(r)["version"])
		if version == "" {
			version = APIVERSION
		}
		writeCorsHeaders(w, r)

		if version.GreaterThan(APIVERSION) {
			http.Error(w, fmt.Errorf("client and server don't have same version (client : %s, server: %s)", version, APIVERSION).Error(), http.StatusNotFound)
			return
		}

		if err := handlerFunc(eng, version, w, r, mux.Vars(r)); err != nil {
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
