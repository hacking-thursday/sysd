// +build linux

package sysfs

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/sysfs", handler)
}

type TreeNode map[string]interface{}

func handler(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	target_path := "/sys/class"
	result := TreeNode{}
	var out []byte

	result["class"] = TreeNode{}

	d_pathes, err := ioutil.ReadDir(target_path)
	if err == nil {
		for i := 0; i < len(d_pathes); i++ {
			target_path2 := path.Join(target_path, d_pathes[i].Name())
			d_pathes2, err2 := ioutil.ReadDir(target_path2)

			if err2 == nil {
				for j := 0; j < len(d_pathes2); j++ {
					target_path3 := path.Join(target_path2, d_pathes2[j].Name())
					target_link3, _ := os.Readlink(target_path3)
					target_link3, _ = filepath.EvalSymlinks(target_path3)
					row := [2]string{target_path3, target_link3}

					key1 := d_pathes[i].Name()
					key2 := d_pathes2[j].Name()

					if _, ok := result["class"].(TreeNode)[key1]; !ok {
						result["class"].(TreeNode)[key1] = TreeNode{}
					}

					if _, ok := result["class"].(TreeNode)[key1].(TreeNode)[key2]; !ok {
						result["class"].(TreeNode)[key1].(TreeNode)[key2] = row
					}

					//result = append(result, row)
				}
			}
		}
	}

	subdir := "dev"
	target_path = "/sys/" + subdir
	result[subdir] = TreeNode{}

	d_pathes, err = ioutil.ReadDir(target_path)
	if err == nil {
		for i := 0; i < len(d_pathes); i++ {
			target_path2 := path.Join(target_path, d_pathes[i].Name())
			d_pathes2, err2 := ioutil.ReadDir(target_path2)

			if err2 == nil {
				for j := 0; j < len(d_pathes2); j++ {
					target_path3 := path.Join(target_path2, d_pathes2[j].Name())
					target_link3, _ := os.Readlink(target_path3)
					target_link3, _ = filepath.EvalSymlinks(target_path3)
					row := [2]string{target_path3, target_link3}

					key1 := d_pathes[i].Name()
					key2 := d_pathes2[j].Name()

					if _, ok := result[subdir].(TreeNode)[key1]; !ok {
						result[subdir].(TreeNode)[key1] = TreeNode{}
					}

					if _, ok := result[subdir].(TreeNode)[key1].(TreeNode)[key2]; !ok {
						result[subdir].(TreeNode)[key1].(TreeNode)[key2] = row
					}
				}
			}
		}
	}

	subdir = "bus"
	target_path = "/sys/" + subdir
	result[subdir] = TreeNode{}

	d_pathes, err = ioutil.ReadDir(target_path)
	if err == nil {
		for i := 0; i < len(d_pathes); i++ {
			key1 := d_pathes[i].Name()

			devices, _ := ioutil.ReadDir(path.Join(target_path, key1, "devices"))
			drivers, _ := ioutil.ReadDir(path.Join(target_path, key1, "drivers"))

			if _, ok := result[subdir].(TreeNode)[key1]; !ok {
				result[subdir].(TreeNode)[key1] = TreeNode{}
			}

			if _, ok := result[subdir].(TreeNode)[key1].(TreeNode)["drivers"]; !ok {
				result[subdir].(TreeNode)[key1].(TreeNode)["drivers"] = TreeNode{}
			}

			if _, ok := result[subdir].(TreeNode)[key1].(TreeNode)["devices"]; !ok {
				result[subdir].(TreeNode)[key1].(TreeNode)["devices"] = TreeNode{}
			}

			for j := 0; j < len(drivers); j++ {
				key2 := drivers[j].Name()
				result[subdir].(TreeNode)[key1].(TreeNode)["drivers"].(TreeNode)[key2] = TreeNode{}
			}

			for j := 0; j < len(devices); j++ {
				key2 := devices[j].Name()
				result[subdir].(TreeNode)[key1].(TreeNode)["devices"].(TreeNode)[key2] = TreeNode{}
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
