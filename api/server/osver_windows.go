package server

import (
	"net/http"
	"runtime"
	"syscall"
)

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

	_, err = w.Write(out)
	return
}
