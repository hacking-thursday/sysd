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

type Battery struct {
    Name string
    Capacity int
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

func get_battery_detail(name string) (Battery) {
    battery := Battery{}
    battery.Name = name
    var myregexp = regexp.MustCompile(`(\w+)=(\w+)`)

    path := "/sys/class/power_supply/" + name + "/uevent"

    ra, _ := ioutil.ReadFile(path)
    list := strings.Split(strings.TrimSpace(string(ra)),"\n")

    for _, value := range list {
        list_a := myregexp.FindStringSubmatch(strings.TrimSpace(value))

        if(list_a[1] == "POWER_SUPPLY_CAPACITY") {
            battery.Capacity, _ = strconv.Atoi(list_a[2])
        }
    }

    return battery
}

func get_batterys(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var output []byte
    var batterys []Battery

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
