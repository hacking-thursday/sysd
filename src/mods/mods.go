package mods

import (
	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/docker/utils"
	"net/http"
	"strings"
)

var (
	Modules = map[string]map[string]HttpApiFunc{
		"GET":     {},
		"POST":    {},
		"DELETE":  {},
		"OPTIONS": {},
	}
)

func init() {
	log.Debugf("pkg mods init()")
}

type HttpApiFunc func(eng *engine.Engine, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) error

func Register(method string, route string, fct HttpApiFunc) error {
	if _, exists := Modules[method][route]; exists {
		log.Debugf("HttpApiFunc already registered %s::%s", method, route)
		return nil
	}

	Modules[method][route] = fct

	return nil
}

//If we don't do this, POST method without Content-type (even with empty body) will fail
func ParseForm(r *http.Request) error {
	if r == nil {
		return nil
	}
	if err := r.ParseForm(); err != nil && !strings.HasPrefix(err.Error(), "mime:") {
		return err
	}
	return nil
}

func StreamJSON(job *engine.Job, w http.ResponseWriter, flush bool) {
	w.Header().Set("Content-Type", "application/json")
	if job.GetenvBool("lineDelim") {
		w.Header().Set("Content-Type", "application/x-json-stream")
	}

	if flush {
		job.Stdout.Add(utils.NewWriteFlusher(w))
	} else {
		job.Stdout.Add(w)
	}
}
