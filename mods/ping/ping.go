package ping

import (
	"github.com/docker/docker/pkg/log"
	"github.com/hacking-thursday/sysd/mods"
	"net/http"
	"strings"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/ping", ping)
}

func ping(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	if _, err = w.Write([]byte("pong")); err != nil {
		httpError(w, err)
		return
	}
	return
}

func httpError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	// FIXME: this is brittle and should not be necessary.
	// If we need to differentiate between different possible error types, we should
	// create appropriate error types with clearly defined meaning.
	if strings.Contains(err.Error(), "No such") {
		statusCode = http.StatusNotFound
	} else if strings.Contains(err.Error(), "Bad parameter") {
		statusCode = http.StatusBadRequest
	} else if strings.Contains(err.Error(), "Conflict") {
		statusCode = http.StatusConflict
	} else if strings.Contains(err.Error(), "Impossible") {
		statusCode = http.StatusNotAcceptable
	} else if strings.Contains(err.Error(), "Wrong login/password") {
		statusCode = http.StatusUnauthorized
	} else if strings.Contains(err.Error(), "hasn't been activated") {
		statusCode = http.StatusForbidden
	}

	if err != nil {
		log.Errorf("HTTP Error: statusCode=%d %s", statusCode, err.Error())
		http.Error(w, err.Error(), statusCode)
	}
}
