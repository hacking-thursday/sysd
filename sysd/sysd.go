package main

import (
	"builtins"
	"daemon"
	"github.com/docker/docker/engine"
	"log"
)

func main() {
	host := "127.0.0.1:4000"

	eng := engine.New()
	if err := builtins.Register(eng); err != nil {
		log.Fatal(err)
	}

	go func() {
		daemonCfg := &daemon.Config{}
		daemonCfg.InstallFlags()
		daemonCfg.Pidfile = "/tmp/sysd.pid"
		daemonCfg.Root = "/tmp"

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

	job := eng.Job("serveapi", "tcp://"+host)
	if err := job.Run(); err != nil {
		log.Fatal(err)
	}
}
