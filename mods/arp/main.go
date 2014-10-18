// +build linux

package arp

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/arp", handler)
}

func handler(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	target_path := "/proc/net/arp"
	result := []map[string]string{}
	header := []string{}
	var out []byte

	f, err := os.Open(target_path)
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		i += 1
		line := scanner.Text()
		fields := strings.Fields(line)

		// 針對第一行處理
		if i == 1 {
			for _, v := range fields {
				header = append(header, v)
			}
			continue
		}

		row := map[string]string{}
		for ii := 0; ii < len(fields); ii++ {
			key := header[ii]
			val := fields[ii]
                        row[key] = val
		}
                result = append( result, row )
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
