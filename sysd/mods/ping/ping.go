package ping

import (
	log "github.com/sirupsen/logrus"
	"sysd/mods"
	"net/http"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/ping", ping)
}

func ping(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	if _, err = w.Write([]byte("pong")); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
