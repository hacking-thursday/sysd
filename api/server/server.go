package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/docker/docker/pkg/log"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/gorilla/mux"
	"github.com/tsaikd/KDGoLib/env"
)

type HttpApiFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) error

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

		if err := handlerFunc(w, r, mux.Vars(r)); err != nil {
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

func marshal(r *http.Request, v interface{}) (b []byte, err error) {
	var query = r.URL.Query()

	if query.Get("pretty") == "1" {
		b, err = json.MarshalIndent(v, "", "\t")
	} else {
		b, err = json.Marshal(v)
	}

	return
}
