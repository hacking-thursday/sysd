// +build docker

package main

import (
	"github.com/docker/docker/daemon"
	"github.com/docker/docker/engine"
	"github.com/hacking-thursday/sysd/builtins"
	"log"
)

func runDaemonByDocker() {
	host := "tcp://127.0.0.1:8"

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
