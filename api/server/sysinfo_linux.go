package server

import (
	"net/http"
	"syscall"
)

func sysinfo(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var (
		out []byte

		si = syscall.Sysinfo_t{}
	)

	if err = syscall.Sysinfo(&si); err != nil {
		httpError(w, err)
		return
	}

	if out, err = marshal(r, si); err != nil {
		httpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		httpError(w, err)
		return
	}

	return
}
