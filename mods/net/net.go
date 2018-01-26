// +build docker

package net

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
	"net"
	"net/http"
)

type Iface struct {
	IP []string
}

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/net", ifconfig)
}

func ifconfig(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	eng := eng_ifce.(*engine.Engine)

	eng.Register("ifconfig", func(job *engine.Job) engine.Status {
		outs := engine.NewTable("", 0)
		out := &engine.Env{}

		ifaces, err := net.Interfaces()
		if err != nil {
			return job.Error(err)
		}
		for _, i := range ifaces {
			iface := Iface{}
			ips := make([]string, 0)

			addrs, err := i.Addrs()
			if err != nil {
				fmt.Print(err)
			}
			for _, a := range addrs {
				ips = append(ips, a.String())
			}

			iface = Iface{ips}
			out.SetJson(i.Name, iface)
		}

		outs.Add(out)

		if _, err := outs.WriteListTo(job.Stdout); err != nil {
			return job.Error(err)
		}

		return engine.StatusOK
	})

	var job = eng.Job("ifconfig", "the_arg0")
	mods.StreamJSON(job, w, false)

	if err := job.Run(); err != nil {
		return err
	}
	return nil
}
