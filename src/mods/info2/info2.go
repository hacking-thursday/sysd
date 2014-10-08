package info2

import (
	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/docker/utils"
	"mods"
	"net/http"
	"os/exec"
	"strings"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/info2", getInfo2)
}

//If we don't do this, POST method without Content-type (even with empty body) will fail
func parseForm(r *http.Request) error {
	if r == nil {
		return nil
	}
	if err := r.ParseForm(); err != nil && !strings.HasPrefix(err.Error(), "mime:") {
		return err
	}
	return nil
}

func streamJSON(job *engine.Job, w http.ResponseWriter, flush bool) {
	w.Header().Set("Content-Type", "application/json")
	if job.GetenvBool("lineDelim") {
		w.Header().Set("Content-Type", "application/x-json-stream")
	}

	if flush {
		job.Stdout.Add(utils.NewWriteFlusher(w))
	} else {
		job.Stdout.Add(w)
	}
}

func getInfo2(eng *engine.Engine, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	eng.Register("my_cmd", func(job *engine.Job) engine.Status {
		outs := engine.NewTable("", 0)

		out := &engine.Env{}

		// 取用命令列參數
		out.Set("arg0", job.Args[0])

		// 取用 Request Header
		out.Set("UA", r.Header.Get("User-Agent"))

		// 取用 Request 的 GET/POST 變數
		if err := parseForm(r); err == nil {
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
	streamJSON(job, w, false)

	if err := job.Run(); err != nil {
		return err
	}
	return nil
}
