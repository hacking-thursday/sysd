// +build linux

package cpuinfo

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/cpuinfo", handler)
}

func handler(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	target_path := "/proc/cpuinfo"
	result := []map[string]interface{}{}
	var out []byte

	f, err := os.Open(target_path)
	scanner := bufio.NewScanner(f)
	i := 0
	row := map[string]interface{}{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			result = append(result, row)
                        row = map[string]interface{}{}

			i += 1
		}

		ary := strings.Split(line, ":")
		if len(ary) == 2 {
			key := strings.TrimSpace(ary[0])
			val := strings.TrimSpace(ary[1])
			row[key] = val

            // Additionally parse flags into array 
			if key == "flags" {
				val_ary := strings.Split(val, " ")
				row[key] = val_ary
			}
		}
	}

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
