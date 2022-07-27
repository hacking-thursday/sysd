package memstats

import (
	log "github.com/sirupsen/logrus"
	"github.com/hacking-thursday/sysd/mods"
	"net/http"
	"runtime"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/memstats", memstats)
}

func memstats(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
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
