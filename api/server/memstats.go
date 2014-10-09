package server

import (
	"net/http"
	"runtime"
)

func memstats(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var (
		m   runtime.MemStats
		out []byte
	)

	runtime.ReadMemStats(&m)
	if out, err = marshal(r, m); err != nil {
		httpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		httpError(w, err)
		return
	}
	return
}
