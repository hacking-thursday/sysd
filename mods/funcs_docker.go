// +build docker

package mods

import (
	"github.com/docker/docker/engine"
	"github.com/docker/docker/utils"
	"net/http"
)

func StreamJSON(job *engine.Job, w http.ResponseWriter, flush bool) {
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
