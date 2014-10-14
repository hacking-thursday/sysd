package main

import (
	"github.com/docker/docker/engine"
	"github.com/hacking-thursday/sysd/builtins"
	"github.com/docker/docker/daemon"
	"log"

	flag "github.com/docker/docker/pkg/mflag"

	apiserver "github.com/hacking-thursday/sysd/api/server2"
)

func runDaemonByDocker() {
	host := "tcp://127.0.0.1:8080"

	eng := engine.New()
	if err := builtins.Register(eng); err != nil {
		log.Fatal(err)
	}

	go func() {
		daemonCfg := &daemon.Config{}
		daemonCfg.InstallFlags()
		daemonCfg.Pidfile = "/tmp/sysd.pid"
		daemonCfg.Root = "/tmp"
		daemonCfg.BridgeIface= "none"

		d, err := daemon.NewDaemon(daemonCfg, eng)
		if err != nil {
			log.Fatal(err)
		}
		if err := d.Install(eng); err != nil {
			log.Fatal(err)
		}

		if err := eng.Job("acceptconnections").Run(); err != nil {
			log.Fatal(err)
		}
	}()

	job := eng.Job("serveapi", host)
	if err := job.Run(); err != nil {
		log.Fatal(err)
	}
}

func runDaemonByKdTsai() {
	apiserver.ListenAndServe()
}

func main() {
	flag.Parse()

	if *flBackend == "minimal" {
		runDaemonByKdTsai()
	} else if *flBackend == "docker" {
		runDaemonByDocker()
	} else {
		runDaemonByKdTsai()
	}
}
