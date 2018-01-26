// +build docker

package info2

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
	"net/http"
	"os/exec"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/info2", getInfo2)
}

func getInfo2(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	eng := eng_ifce.(*engine.Engine)

	eng.Register("my_cmd", func(job *engine.Job) engine.Status {
		outs := engine.NewTable("", 0)

		out := &engine.Env{}

		// 取用命令列參數
		out.Set("arg0", job.Args[0])

		// 取用 Request Header
		out.Set("UA", r.Header.Get("User-Agent"))

		// 取用 Request 的 GET/POST 變數
		if err := mods.ParseForm(r); err == nil {
			out.Set("var0", r.Form.Get("var0"))
		}

		// 取用系統指令輸出
		output, err := exec.Command("uname", "-a").Output()
		if err == nil {
			out.Set("uname", string(output))
		}

		outs.Add(out)

		if _, err := outs.WriteListTo(job.Stdout); err != nil {
			return job.Error(err)
		}

		return engine.StatusOK
	})

	var job = eng.Job("my_cmd", "the_arg0")
	mods.StreamJSON(job, w, false)

	if err := job.Run(); err != nil {
		return err
	}
	return nil
}
