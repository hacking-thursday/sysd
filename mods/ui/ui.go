package ui

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/docker/docker/pkg/version"
	"github.com/tsaikd/KDGoLib/env"
	"github.com/tsaikd/KDGoLib/futil"

	"github.com/hacking-thursday/sysd/mods"
)

var (
	flUiDir = flag.String(
		[]string{"-SYSD_UI_DIR"},
		env.GetString("SYSD_UI_DIR", ""),
		"SYSD UI Directory",
	)
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/ui/*", handler_ui)

	if !futil.IsExist(*flUiDir + "/index.html") {
		tryPaths := []string{"files", "../mods/ui/files", "/usr/share/sysd/webui"}
		for _, path := range tryPaths {
			if futil.IsExist(path + "/index.html") {
				*flUiDir = path
				break
			}

		}

		if !futil.IsExist(*flUiDir + "/index.html") {
			log.Warnf("Incorrent UI directory: %v", *flUiDir)
		}
	}
}

func handler_ui(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	if len(r.URL.Path) < 4 {
		http.Redirect(w, r, "/ui/", http.StatusMovedPermanently)
		return
	}

	doc_root := strings.TrimSuffix(*flUiDir, "/")

	subpath := r.URL.Path[4:]
	if subpath == "" {
		subpath = "index.html"
	}

	f_path := doc_root + "/" + subpath

	log.Debug("handler_ui: ", f_path)
	http.ServeFile(w, r, f_path)

	return
}
