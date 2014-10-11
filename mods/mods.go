package mods

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/pkg/log"
	"github.com/gorilla/mux"
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

func makeHttpHandler(localMethod string, localRoute string, handlerFunc HttpApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log the request
		log.Debugf("Calling %s %s", localMethod, localRoute)
		log.Infof("%s %s", r.Method, r.RequestURI)

		if err := handlerFunc(w, r, mux.Vars(r)); err != nil {
			log.Errorf("Handler for %s %s returned error: %s", localMethod, localRoute, err)
			HttpError(w, err)
		}
	}
}

func HttpError(w http.ResponseWriter, err error) {
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

func Marshal(r *http.Request, v interface{}) (b []byte, err error) {
	var query = r.URL.Query()

	if query.Get("pretty") == "1" {
		b, err = json.MarshalIndent(v, "", "\t")
	} else {
		b, err = json.Marshal(v)
	}

	return
}
