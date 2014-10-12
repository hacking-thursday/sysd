package server

import (
	"net/http"

	"github.com/docker/docker/pkg/version"

	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	mods.Register("GET", "/ping", ping)
}

func ping(engine interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	if _, err = w.Write([]byte("pong")); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
