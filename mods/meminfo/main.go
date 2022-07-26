// +build linux

package meminfo

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/meminfo", handler)
}

func handler(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	target_path := "/proc/meminfo"
	result := []map[string]string{}
	var out []byte

	f, err := os.Open(target_path)
	scanner := bufio.NewScanner(f)
	i := 0
	row := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			result = append(result, row)
			row = map[string]string{}

			i += 1
		}

		ary := strings.Split(line, ":")
		if len(ary) == 2 {
			key := strings.TrimSpace(ary[0])
			val := strings.TrimSpace(ary[1])
			row[key] = val
		}
	}
	result = append(result, row)

	if out, err = mods.Marshal(r, result); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
