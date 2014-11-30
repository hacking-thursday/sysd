// +build linux

package process

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/process/resource", handler)
}

type TreeNode map[string]interface{}

func handler(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
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
                        f.Close()

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
						result["processes"].(TreeNode)[pid].(TreeNode)["cmdline"] = cmdline_txt
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
