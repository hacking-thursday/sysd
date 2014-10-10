package server2

import (
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	_, err = w.Write([]byte("pong"))
	return
}
