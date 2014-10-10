package server2

import (
	"net/http"

	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/version"
)

func ping(eng *engine.Engine, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	_, err = w.Write([]byte("pong"))
	return
}
