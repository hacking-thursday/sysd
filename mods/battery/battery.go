package battery

import (
    "github.com/docker/docker/pkg/log"
    "github.com/docker/docker/pkg/version"
    "github.com/hacking-thursday/sysd/mods"
    "net/http"
    "io/ioutil"
    "regexp"
    "strings"
    "strconv"
)

type bat struct {
    Capacity int
}

var myregexp = regexp.MustCompile(`(\w+)=(\w+)`)

func init() {
    log.Debugf("Initializing module...")
    mods.Register("GET", "/battery", get_battery)
}

func get_battery(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
    ra, _ := ioutil.ReadFile("/sys/class/power_supply/BAT0/uevent")
    list := strings.Split(strings.TrimSpace(string(ra)),"\n")
	var output []byte
    bat_info := bat{}

    for _, value := range list {
        list_a := myregexp.FindStringSubmatch(strings.TrimSpace(value))

        if(list_a[1] == "POWER_SUPPLY_CAPACITY") {
            bat_info.Capacity, _ = strconv.Atoi(list_a[2])
        }
    }

	if output, err = mods.Marshal(r, bat_info); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(output); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
