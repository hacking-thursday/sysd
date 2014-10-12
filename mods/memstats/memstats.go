package server

import (
	"net/http"
	"runtime"

	"github.com/docker/docker/pkg/version"

	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	mods.Register("GET", "/memstats", memstats)
}

func memstats(engine interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var (
		m   runtime.MemStats
		out []byte
	)

	runtime.ReadMemStats(&m)
	if out, err = mods.Marshal(r, m); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
