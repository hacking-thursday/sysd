package server

import (
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	if _, err = w.Write([]byte("pong")); err != nil {
		httpError(w, err)
		return
	}
	return
}
