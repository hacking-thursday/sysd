package ui

import (
	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
	"net/http"
	"os"
)

// 先用很醜的方法，暫時先支援 3 層目錄 XD
func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/ui", handler_ui)
	mods.Register("GET", "/ui/", handler_ui)
	mods.Register("GET", "/ui/{.*}", handler_ui)
	mods.Register("GET", "/ui/{.*}/", handler_ui)
	mods.Register("GET", "/ui/{.*}/{.*}", handler_ui)
	mods.Register("GET", "/ui/{.*}/{.*}/", handler_ui)
	mods.Register("GET", "/ui/{.*}/{.*}/{.*}", handler_ui)
	mods.Register("GET", "/ui/{.*}/{.*}/{.*}/", handler_ui)
}

func handler_ui(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	doc_root := os.Getenv("SYSD_UI_DIR")

	subpath := r.URL.Path[3:]
	if subpath == "/" || subpath == "" {
		subpath = "index.html"
	}

	f_path := doc_root + "/" + subpath

	log.Debugf("f_path :: " + f_path)
	http.ServeFile(w, r, f_path)

	return
}
