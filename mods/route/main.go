// +build linux

package route

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/route", handler)
	mods.Register("GET", "/network/route", handler)
}

func hex_to_ip(input string) string {
	a, _ := hex.DecodeString(input)
	return fmt.Sprintf("%v.%v.%v.%v", a[3], a[2], a[1], a[0])
}

func handler(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	target_path := "/proc/net/route"
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
			if key == "Destination" || key == "Gateway" || key == "Mask" {
				val = hex_to_ip(val)
			}
			row[key] = val
		}
		result = append(result, row)
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
