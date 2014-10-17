package sample

import (
	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
	"net/http"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/sample", handler_sample)
}

func handler_sample(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
        w.Write([]byte("hello world!! sample"))

	return
}
