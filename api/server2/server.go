package server2

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/docker/docker/pkg/log"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/gorilla/mux"
	"github.com/tsaikd/KDGoLib/env"

	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/version"

	"github.com/hacking-thursday/sysd/mods"
)

type HttpApiFunc func(eng *engine.Engine, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) error

var (
	flApiAddr = flag.String(
		[]string{"-SYSD_API_ADDR"},
		env.GetString("SYSD_API_ADDR", "tcp://0.0.0.0:8080"),
		"Sysd API Server Listen Address",
	)
	flApiPrefix = flag.String(
		[]string{"-SYSD_API_PREFIX"},
		env.GetString("SYSD_API_PREFIX", ""),
		"Sysd API Server URL Prefix",
	)
)

func ListenAndServe() (err error) {
	var (
		l       net.Listener
		r       *mux.Router
		apiaddr = *flApiAddr
		prefix  = *flApiPrefix
		proto   string
		addr    string
	)

	if r, err = createRouter(); err != nil {
		return
	}

	if proto, addr, err = parseAddr(apiaddr); err != nil {
		return
	}

	if l, err = net.Listen(proto, addr); err != nil {
		return
	}

	httpSrv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Infof("Listening for HTTP on %s (%s)", addr, prefix)
	err = httpSrv.Serve(l)

	return
}

func createRouter() (r *mux.Router, err error) {
	var (
		prefix = *flApiPrefix
	)
	r = mux.NewRouter()

	m := map[string]map[string]HttpApiFunc{
		"GET": {
			"/ping": ping,
		},
	}

	// beg 載入並註冊自定義的處理函式模組
	for method, routes := range m {
		routes2 := mods.Modules[method]
		for route, fct := range routes2 {
			if _, exists := routes[route]; exists {
				continue
			}
			m[method][route] = HttpApiFunc(fct)
		}
	}
	// end

	for method, routes := range m {
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

func parseAddr(apiaddr string) (proto string, addr string, err error) {
	seps := strings.Split(apiaddr, "://")
	if len(seps) < 2 {
		err = fmt.Errorf("Invalid API Address format", apiaddr)
		return
	}

	proto = seps[0]
	addr = seps[1]

	if proto != "tcp" {
		err = fmt.Errorf("API Address support only tcp now")
		proto = ""
		addr = ""
		return
	}

	return
}

func makeHttpHandler(localMethod string, localRoute string, handlerFunc HttpApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log the request
		log.Debugf("Calling %s %s", localMethod, localRoute)
		log.Infof("%s %s", r.Method, r.RequestURI)

		if err := handlerFunc(nil, "", w, r, mux.Vars(r)); err != nil {
			log.Errorf("Handler for %s %s returned error: %s", localMethod, localRoute, err)
			httpError(w, err)
		}
	}
}

func httpError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	// FIXME: this is brittle and should not be necessary.
	// If we need to differentiate between different possible error types, we should
	// create appropriate error types with clearly defined meaning.
	if strings.Contains(err.Error(), "No such") {
		statusCode = http.StatusNotFound
	} else if strings.Contains(err.Error(), "Bad parameter") {
		statusCode = http.StatusBadRequest
	} else if strings.Contains(err.Error(), "Conflict") {
		statusCode = http.StatusConflict
	} else if strings.Contains(err.Error(), "Impossible") {
		statusCode = http.StatusNotAcceptable
	} else if strings.Contains(err.Error(), "Wrong login/password") {
		statusCode = http.StatusUnauthorized
	} else if strings.Contains(err.Error(), "hasn't been activated") {
		statusCode = http.StatusForbidden
	}

	if err != nil {
		log.Errorf("HTTP Error: statusCode=%d %s", statusCode, err.Error())
		http.Error(w, err.Error(), statusCode)
	}
}
