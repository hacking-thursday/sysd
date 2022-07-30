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
	"sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/network/socket", handler)
}

func hex_to_ip(input string) string {
	a, _ := hex.DecodeString(input)
	return fmt.Sprintf("%v.%v.%v.%v", a[3], a[2], a[1], a[0])
}

func handler(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	result := []map[string]string{}
	header := []string{}
	var out []byte

	// unix
	target_path := "/proc/net/unix"
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
				header = append(header, strings.ToLower(v))
			}
			continue
		}

		row := map[string]string{}
		row["socket_type"] = "unix"
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

	// tcp
	target_path = "/proc/net/tcp"
	header = []string{}
	f, err = os.Open(target_path)
	scanner = bufio.NewScanner(f)
	i = 0
	for scanner.Scan() {
		i += 1
		line := scanner.Text()
		fields := strings.Fields(line)

		// 針對第一行處理
		if i == 1 {
			for _, v := range fields {
				header = append(header, strings.ToLower(v))
			}
			continue
		}

		row := map[string]string{}
		row["socket_type"] = "tcp"

		{
			key := header[0]
			val := fields[0]
			row[key] = val
		}

		{
			key := header[1]
			val := fields[1]
			row[key] = val
		}

		{
			key := header[2]
			val := fields[2]
			row[key] = val
		}

		{
			key := header[3]
			val := fields[3]
			row[key] = val
		}

		{
			key := header[4]
			val := strings.Split(fields[4], ":")[0]
			row[key] = val
		}

		{
			key := header[5]
			val := strings.Split(fields[4], ":")[1]
			row[key] = val
		}

		{
			key := header[6]
			val := strings.Split(fields[5], ":")[0]
			row[key] = val
		}

		{
			key := header[7]
			val := strings.Split(fields[5], ":")[1]
			row[key] = val
		}

		{
			key := header[8]
			val := fields[6]
			row[key] = val
		}

		{
			key := header[9]
			val := fields[7]
			row[key] = val
		}

		{
			key := header[10]
			val := fields[8]
			row[key] = val
		}

		{
			key := header[11]
			val := fields[9]
			row[key] = val
		}

		result = append(result, row)
	} // tcp

	target_path = "/proc/net/tcp6"
	header = []string{}
	f, err = os.Open(target_path)
	scanner = bufio.NewScanner(f)
	i = 0
	for scanner.Scan() {
		i += 1
		line := scanner.Text()
		fields := strings.Fields(line)

		// 針對第一行處理
		if i == 1 {
			for _, v := range fields {
				header = append(header, strings.ToLower(v))
			}
			continue
		}

		row := map[string]string{}
		row["socket_type"] = "tcp"

		{
			key := header[0]
			val := fields[0]
			row[key] = val
		}

		{
			key := header[1]
			val := fields[1]
			row[key] = val
		}

		{
			key := header[2]
			val := fields[2]
			row[key] = val
		}

		{
			key := header[3]
			val := fields[3]
			row[key] = val
		}

		{
			key := header[4]
			val := strings.Split(fields[4], ":")[0]
			row[key] = val
		}

		{
			key := header[5]
			val := strings.Split(fields[4], ":")[1]
			row[key] = val
		}

		{
			key := header[6]
			val := strings.Split(fields[5], ":")[0]
			row[key] = val
		}

		{
			key := header[7]
			val := strings.Split(fields[5], ":")[1]
			row[key] = val
		}

		{
			key := header[8]
			val := fields[6]
			row[key] = val
		}

		{
			key := header[9]
			val := fields[7]
			row[key] = val
		}

		{
			key := header[10]
			val := fields[8]
			row[key] = val
		}

		{
			key := header[11]
			val := fields[9]
			row[key] = val
		}

		result = append(result, row)
	}

	// ucp
	target_path = "/proc/net/udp"
	header = []string{}
	f, err = os.Open(target_path)
	scanner = bufio.NewScanner(f)
	i = 0
	for scanner.Scan() {
		i += 1
		line := scanner.Text()
		fields := strings.Fields(line)

		// 針對第一行處理
		if i == 1 {
			for _, v := range fields {
				header = append(header, strings.ToLower(v))
			}
			continue
		}

		row := map[string]string{}
		row["socket_type"] = "udp"

		{
			key := header[0]
			val := fields[0]
			row[key] = val
		}

		{
			key := header[1]
			val := fields[1]
			row[key] = val
		}

		{
			key := header[2]
			val := fields[2]
			row[key] = val
		}

		{
			key := header[3]
			val := fields[3]
			row[key] = val
		}

		{
			key := header[4]
			val := strings.Split(fields[4], ":")[0]
			row[key] = val
		}

		{
			key := header[5]
			val := strings.Split(fields[4], ":")[1]
			row[key] = val
		}

		{
			key := header[6]
			val := strings.Split(fields[5], ":")[0]
			row[key] = val
		}

		{
			key := header[7]
			val := strings.Split(fields[5], ":")[1]
			row[key] = val
		}

		{
			key := header[8]
			val := fields[6]
			row[key] = val
		}

		{
			key := header[9]
			val := fields[7]
			row[key] = val
		}

		{
			key := header[10]
			val := fields[8]
			row[key] = val
		}

		{
			key := header[11]
			val := fields[9]
			row[key] = val
		}

		result = append(result, row)
	}

	// udp
	target_path = "/proc/net/udp6"
	header = []string{}
	f, err = os.Open(target_path)
	scanner = bufio.NewScanner(f)
	i = 0
	for scanner.Scan() {
		i += 1
		line := scanner.Text()
		fields := strings.Fields(line)

		// 針對第一行處理
		if i == 1 {
			for _, v := range fields {
				header = append(header, strings.ToLower(v))
			}
			continue
		}

		row := map[string]string{}
		row["socket_type"] = "udp"

		{
			key := header[0]
			val := fields[0]
			row[key] = val
		}

		{
			key := header[1]
			val := fields[1]
			row[key] = val
		}

		{
			key := header[2]
			val := fields[2]
			row[key] = val
		}

		{
			key := header[3]
			val := fields[3]
			row[key] = val
		}

		{
			key := header[4]
			val := strings.Split(fields[4], ":")[0]
			row[key] = val
		}

		{
			key := header[5]
			val := strings.Split(fields[4], ":")[1]
			row[key] = val
		}

		{
			key := header[6]
			val := strings.Split(fields[5], ":")[0]
			row[key] = val
		}

		{
			key := header[7]
			val := strings.Split(fields[5], ":")[1]
			row[key] = val
		}

		{
			key := header[8]
			val := fields[6]
			row[key] = val
		}

		{
			key := header[9]
			val := fields[7]
			row[key] = val
		}

		{
			key := header[10]
			val := fields[8]
			row[key] = val
		}

		{
			key := header[11]
			val := fields[9]
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
