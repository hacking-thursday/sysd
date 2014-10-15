package server2

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/docker/docker/pkg/log"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/gorilla/mux"
	"github.com/tsaikd/KDGoLib/env"

	"github.com/hacking-thursday/sysd/mods"
	_ "github.com/hacking-thursday/sysd/mods/loader"
)

var (
	flApiAddr = flag.String(
		[]string{"-SYSD_API_ADDR"},
		env.GetString("SYSD_API_ADDR", "tcp://0.0.0.0:8080"),
		"Sysd API Server Listen Address",
	)
)

func ListenAndServe(eng interface{}) (err error) {
	var (
		l       net.Listener
		r       *mux.Router
		apiaddr = *flApiAddr
		proto   string
		addr    string
	)

	if r, err = mods.CreateRouter(eng); err != nil {
		return
	}

	if proto, addr, err = parseAddr(apiaddr); err != nil {
		return
	}

	if l, err = net.Listen(proto, addr); err != nil {
		return
	}

	httpSrv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Infof("Listening for HTTP on %s", addr)
	err = httpSrv.Serve(l)

	return
}

func parseAddr(apiaddr string) (proto string, addr string, err error) {
	seps := strings.Split(apiaddr, "://")
	if len(seps) < 2 {
		err = fmt.Errorf("Invalid API Address format", apiaddr)
		return
	}

	proto = seps[0]
	addr = seps[1]

	if proto != "tcp" {
		err = fmt.Errorf("API Address support only tcp now")
		proto = ""
		addr = ""
		return
	}

	return
}
