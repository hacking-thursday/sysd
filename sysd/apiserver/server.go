package apiserver

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "flag"
	"github.com/gorilla/mux"
	"github.com/tsaikd/KDGoLib/env"

	"sysd/mods"
	_ "sysd/mods/loader"
)

var (
	flApiAddr = flag.String(
		"SYSD_API_ADDR",
		env.GetString("SYSD_API_ADDR", "tcp://0.0.0.0:8"),
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
