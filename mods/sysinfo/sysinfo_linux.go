package sysinfo

import (
	log "github.com/sirupsen/logrus"
	"github.com/hacking-thursday/sysd/mods"
	"net/http"
	"syscall"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/sysinfo", sysinfo)
}

func sysinfo(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
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
