package battery

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"sysd/mods"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func MakeFirstUpperCase(s string) string {

	if len(s) < 2 {
		return strings.ToUpper(s)
	}

	bts := []byte(strings.ToLower(s))

	lc := bytes.ToUpper([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/battery", get_batterys)
}

func get_battery_names() []string {
	var batterys []string
	r, _ := regexp.Compile("BAT.*")

	list, _ := ioutil.ReadDir("/sys/class/power_supply/")

	for _, dir := range list {
		if r.MatchString(dir.Name()) {
			batterys = append(batterys, dir.Name())
		}
	}

	return batterys
}

func get_battery_detail(name string) map[string]string {
	var battery = make(map[string]string)
	var myregexp = regexp.MustCompile(`POWER_SUPPLY_(.*)=(.*)`)

	path := "/sys/class/power_supply/" + name + "/uevent"

	ra, _ := ioutil.ReadFile(path)
	list := strings.Split(strings.TrimSpace(string(ra)), "\n")

	for _, value := range list {
		list_a := myregexp.FindStringSubmatch(strings.TrimSpace(value))

		/* make property name "MODEL_NAME" -> "ModelName" */
		var key_name_split = strings.Split(list_a[1], "_")
		var key_name string
		for _, value := range key_name_split {
			key_name = key_name + MakeFirstUpperCase(value)
		}
		/* make property name "MODEL_NAME" -> "ModelName" */

		battery[key_name] = list_a[2]
	}

	return battery
}

func get_batterys(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var output []byte
	var batterys []map[string]string

	for _, name := range get_battery_names() {
		batterys = append(batterys, get_battery_detail(name))
	}

	if output, err = mods.Marshal(r, batterys); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(output); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
