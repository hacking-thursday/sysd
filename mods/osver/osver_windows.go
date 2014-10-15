package osver

import (
	"encoding/json"
	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
	"net/http"
	"runtime"
	"strings"
	"syscall"
)

func init() {
	//log.Debugf("Initializing module...")
	mods.Register("GET", "/osver", osver)
}

func osver(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var (
		ver uint32
		out []byte

		ov struct {
			OS    string
			Arch  string
			Major byte
			Minor uint8
			Build uint16
		}
	)

	ver, err = syscall.GetVersion()

	ov.OS = runtime.GOOS
	ov.Arch = runtime.GOARCH
	ov.Major = byte(ver)
	ov.Minor = uint8(ver >> 8)
	ov.Build = uint16(ver >> 16)

	if out, err = mods.Marshal(r, ov); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
