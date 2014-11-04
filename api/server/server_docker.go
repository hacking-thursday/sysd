// +build docker

package server

import (
	"crypto/tls"
	"crypto/x509"
	"expvar"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/libcontainer/user"
	"github.com/gorilla/mux"

	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/listenbuffer"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/systemd"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/api"
	"github.com/hacking-thursday/sysd/mods"
)

var (
	activationLock chan struct{}
)

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

func getBoolParam(value string) (bool, error) {
	if value == "" {
		return false, nil
	}
	ret, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("Bad parameter")
	}
	return ret, nil
}

func writeCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

func makeHttpHandler(eng *engine.Engine, logging bool, localMethod string, localRoute string, handlerFunc mods.HttpApiFunc, enableCors bool, dockerVersion version.Version) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log the request
		log.Debugf("Calling %s %s", localMethod, localRoute)

		if logging {
			log.Infof("%s %s", r.Method, r.RequestURI)
		}

		if strings.Contains(r.Header.Get("User-Agent"), "Docker-Client/") {
			userAgent := strings.Split(r.Header.Get("User-Agent"), "/")
			if len(userAgent) == 2 && !dockerVersion.Equal(version.Version(userAgent[1])) {
				log.Debugf("Warning: client and server don't have the same version (client: %s, server: %s)", userAgent[1], dockerVersion)
			}
		}
		version := version.Version(mux.Vars(r)["version"])
		if version == "" {
			version = api.APIVERSION
		}
		if enableCors {
			writeCorsHeaders(w, r)
		}

		if version.GreaterThan(api.APIVERSION) {
			http.Error(w, fmt.Errorf("client and server don't have same version (client : %s, server: %s)", version, api.APIVERSION).Error(), http.StatusNotFound)
			return
		}

		if err := handlerFunc(eng, version, w, r, mux.Vars(r)); err != nil {
			log.Errorf("Handler for %s %s returned error: %s", localMethod, localRoute, err)
			httpError(w, err)
		}
	}
}

// Replicated from expvar.go as not public.
func expvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func AttachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/vars", expvarHandler)
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	router.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	router.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

func createRouter(eng *engine.Engine, logging, enableCors bool, dockerVersion string) (*mux.Router, error) {
	r := mux.NewRouter()
	if os.Getenv("DEBUG") != "" {
		AttachProfiler(r)
	}
	m := map[string]map[string]mods.HttpApiFunc{
		"GET":     {},
		"POST":    {},
		"DELETE":  {},
		"OPTIONS": {},
	}

	// beg 載入並註冊自定義的處理函式模組
	for method, routes := range m {
		routes2 := mods.Modules[method]
		for route, fct := range routes2 {
			if _, exists := routes[route]; exists {
				continue
			}
			m[method][route] = fct
		}
	}
	// end

	for method, routes := range m {
		for route, fct := range routes {
			log.Debugf("Registering %s, %s", method, route)
			// NOTE: scope issue, make sure the variables are local and won't be changed
			localRoute := route
			localFct := fct
			localMethod := method

			// build the handler function
			f := makeHttpHandler(eng, logging, localMethod, localRoute, localFct, enableCors, version.Version(dockerVersion))

			// add the new route
			if localRoute == "" {
				r.Methods(localMethod).HandlerFunc(f)
			} else {
				r.Path("/v{version:[0-9.]+}" + localRoute).Methods(localMethod).HandlerFunc(f)
				r.Path(localRoute).Methods(localMethod).HandlerFunc(f)
			}
		}
	}

	return r, nil
}

// ServeRequest processes a single http request to the docker remote api.
// FIXME: refactor this to be part of Server and not require re-creating a new
// router each time. This requires first moving ListenAndServe into Server.
func ServeRequest(eng *engine.Engine, apiversion version.Version, w http.ResponseWriter, req *http.Request) error {
	router, err := createRouter(eng, false, true, "")
	if err != nil {
		return err
	}
	// Insert APIVERSION into the request as a convenience
	req.URL.Path = fmt.Sprintf("/v%s%s", apiversion, req.URL.Path)
	router.ServeHTTP(w, req)
	return nil
}

// ServeFD creates an http.Server and sets it up to serve given a socket activated
// argument.
func ServeFd(addr string, handle http.Handler) error {
	ls, e := systemd.ListenFD(addr)
	if e != nil {
		return e
	}

	chErrors := make(chan error, len(ls))

	// We don't want to start serving on these sockets until the
	// daemon is initialized and installed. Otherwise required handlers
	// won't be ready.
	<-activationLock

	// Since ListenFD will return one or more sockets we have
	// to create a go func to spawn off multiple serves
	for i := range ls {
		listener := ls[i]
		go func() {
			httpSrv := http.Server{Handler: handle}
			chErrors <- httpSrv.Serve(listener)
		}()
	}

	for i := 0; i < len(ls); i++ {
		err := <-chErrors
		if err != nil {
			return err
		}
	}

	return nil
}

func lookupGidByName(nameOrGid string) (int, error) {
	groups, err := user.ParseGroupFilter(func(g *user.Group) bool {
		return g.Name == nameOrGid || strconv.Itoa(g.Gid) == nameOrGid
	})
	if err != nil {
		return -1, err
	}
	if groups != nil && len(groups) > 0 {
		return groups[0].Gid, nil
	}
	return -1, fmt.Errorf("Group %s not found", nameOrGid)
}

func changeGroup(addr string, nameOrGid string) error {
	gid, err := lookupGidByName(nameOrGid)
	if err != nil {
		return err
	}

	log.Debugf("%s group found. gid: %d", nameOrGid, gid)
	return os.Chown(addr, 0, gid)
}

// ListenAndServe sets up the required http.Server and gets it listening for
// each addr passed in and does protocol specific checking.
func ListenAndServe(proto, addr string, job *engine.Job) error {
	var l net.Listener
	r, err := createRouter(job.Eng, job.GetenvBool("Logging"), job.GetenvBool("EnableCors"), job.Getenv("Version"))
	if err != nil {
		return err
	}

	if proto == "fd" {
		return ServeFd(addr, r)
	}

	if proto == "unix" {
		if err := syscall.Unlink(addr); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	var oldmask int
	if proto == "unix" {
		oldmask = syscall.Umask(0777)
	}

	if job.GetenvBool("BufferRequests") {
		l, err = listenbuffer.NewListenBuffer(proto, addr, activationLock)
	} else {
		l, err = net.Listen(proto, addr)
	}

	if proto == "unix" {
		syscall.Umask(oldmask)
	}
	if err != nil {
		return err
	}

	if proto != "unix" && (job.GetenvBool("Tls") || job.GetenvBool("TlsVerify")) {
		tlsCert := job.Getenv("TlsCert")
		tlsKey := job.Getenv("TlsKey")
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			return fmt.Errorf("Couldn't load X509 key pair (%s, %s): %s. Key encrypted?",
				tlsCert, tlsKey, err)
		}
		tlsConfig := &tls.Config{
			NextProtos:   []string{"http/1.1"},
			Certificates: []tls.Certificate{cert},
		}
		if job.GetenvBool("TlsVerify") {
			certPool := x509.NewCertPool()
			file, err := ioutil.ReadFile(job.Getenv("TlsCa"))
			if err != nil {
				return fmt.Errorf("Couldn't read CA certificate: %s", err)
			}
			certPool.AppendCertsFromPEM(file)

			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = certPool
		}
		l = tls.NewListener(l, tlsConfig)
	}

	// Basic error and sanity checking
	switch proto {
	case "tcp":
		if !strings.HasPrefix(addr, "127.0.0.1") && !job.GetenvBool("TlsVerify") {
			log.Infof("/!\\ DON'T BIND ON ANOTHER IP ADDRESS THAN 127.0.0.1 IF YOU DON'T KNOW WHAT YOU'RE DOING /!\\")
		}
	case "unix":
		socketGroup := job.Getenv("SocketGroup")
		if socketGroup != "" {
			if err := changeGroup(addr, socketGroup); err != nil {
				if socketGroup == "docker" {
					// if the user hasn't explicitly specified the group ownership, don't fail on errors.
					log.Debugf("Warning: could not chgrp %s to docker: %s", addr, err.Error())
				} else {
					return err
				}
			}

		}
		if err := os.Chmod(addr, 0660); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid protocol format.")
	}

	httpSrv := http.Server{Addr: addr, Handler: r}
	return httpSrv.Serve(l)
}

// ServeApi loops through all of the protocols sent in to docker and spawns
// off a go routine to setup a serving http.Server for each.
func ServeApi(job *engine.Job) engine.Status {
	if len(job.Args) == 0 {
		return job.Errorf("usage: %s PROTO://ADDR [PROTO://ADDR ...]", job.Name)
	}
	var (
		protoAddrs = job.Args
		chErrors   = make(chan error, len(protoAddrs))
	)
	activationLock = make(chan struct{})

	for _, protoAddr := range protoAddrs {
		protoAddrParts := strings.SplitN(protoAddr, "://", 2)
		if len(protoAddrParts) != 2 {
			return job.Errorf("usage: %s PROTO://ADDR [PROTO://ADDR ...]", job.Name)
		}
		go func() {
			log.Infof("Listening for HTTP on %s (%s)", protoAddrParts[0], protoAddrParts[1])
			chErrors <- ListenAndServe(protoAddrParts[0], protoAddrParts[1], job)
		}()
	}

	for i := 0; i < len(protoAddrs); i++ {
		err := <-chErrors
		if err != nil {
			return job.Error(err)
		}
	}

	return engine.StatusOK
}

func AcceptConnections(job *engine.Job) engine.Status {
	// Tell the init daemon we are accepting requests
	go systemd.SdNotify("READY=1")

	// close the lock so the listeners start accepting connections
	if activationLock != nil {
		close(activationLock)
	}

	return engine.StatusOK
}
