package server

import (
	"net/http"
	"syscall"

	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	mods.Register("GET", "/sysinfo", sysinfo)
}

func sysinfo(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var (
		out []byte

		si = syscall.Sysinfo_t{}
	)

	if err = syscall.Sysinfo(&si); err != nil {
		mods.HttpError(w, err)
		return
	}

	if out, err = mods.Marshal(r, si); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		mods.HttpError(w, err)
		return
	}

	return
}
