package osver

import (
	"net/http"
	"runtime"
	"syscall"
)

func init() {
	//log.Debugf("Initializing module...")
	mods.Register("GET", "/osver", osver)
}

func osver(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
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

	if out, err = marshal(r, ov); err != nil {
		httpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
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

func marshal(r *http.Request, v interface{}) (b []byte, err error) {
	var query = r.URL.Query()

	if query.Get("pretty") == "1" {
		b, err = json.MarshalIndent(v, "", "\t")
	} else {
		b, err = json.Marshal(v)
	}

	return
}
