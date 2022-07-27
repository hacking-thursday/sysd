// +build linux

package process

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/process/resource", handler)
}

type TreeNode map[string]interface{}

func handler(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	target_path := "/proc"
	result := TreeNode{}
	var out []byte

	result["processes"] = TreeNode{}

	d_pathes, err := ioutil.ReadDir(target_path)
	if err == nil {
		for i := 0; i < len(d_pathes); i++ {
			iii, _ := strconv.ParseInt(d_pathes[i].Name(), 0, 64)
			if iii == 0 {
				// 跳過非整數的項目
				continue
			}
			target_path2 := path.Join(target_path, d_pathes[i].Name(), "fd")
			d_pathes2, err2 := ioutil.ReadDir(target_path2)

			cmdline_path := path.Join(target_path, d_pathes[i].Name(), "cmdline")
			f, _ := os.Open(cmdline_path)
			scanner := bufio.NewScanner(f)
			scanner.Scan()
			cmdline_txt := scanner.Text()
			cmdline_ary := strings.Split(cmdline_txt, "\u0000")
			ary_len := len(cmdline_ary)
			if cmdline_ary[ary_len-1] == "" && ary_len >= 2 {
				cmdline_ary = cmdline_ary[0 : ary_len-1]
			}

			f.Close()

			// 讀取 /proc/<pid>/status 的資料
			status_ary := map[string]string{}
			{
				the_path := path.Join(target_path, d_pathes[i].Name(), "status")
				f, _ := os.Open(the_path)
				scanner := bufio.NewScanner(f)
				row := map[string]string{}
				for scanner.Scan() {
					line := scanner.Text()
					// if line == "" {
					// 	result = append(result, row)
					// 	row = map[string]string{}

					// 	i += 1
					// }

					ary := strings.Split(line, ":")
					if len(ary) == 2 {
						key := strings.TrimSpace(ary[0])
						val := strings.TrimSpace(ary[1])
						key = strings.ToLower(key)
						row[key] = val
					}
				}
				status_ary = row
			}

			if err2 == nil {
				for j := 0; j < len(d_pathes2); j++ {
					target_path3 := path.Join(target_path2, d_pathes2[j].Name())
					target_link3, _ := os.Readlink(target_path3)
					row := [2]string{target_path3, target_link3}

					if target_link3 == "" {
						// 跳過空白的項目
						continue
					}

					pid := d_pathes[i].Name()
					fd := d_pathes2[j].Name()

					if _, ok := result["processes"].(TreeNode)[pid]; !ok {
						result["processes"].(TreeNode)[pid] = TreeNode{}
						result["processes"].(TreeNode)[pid].(TreeNode)["fd"] = TreeNode{}
						result["processes"].(TreeNode)[pid].(TreeNode)["cmdline"] = cmdline_ary
						result["processes"].(TreeNode)[pid].(TreeNode)["status"] = status_ary
					}

					if _, ok := result["processes"].(TreeNode)[pid].(TreeNode)["fd"].(TreeNode)[fd]; !ok {
						result["processes"].(TreeNode)[pid].(TreeNode)["fd"].(TreeNode)[fd] = row
					}

					//result = append(result, row)
				}
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
