/*
References :
    Linux Kernel Documentation :: filesystems : proc.txt http://www.mjmwired.net/kernel/Documentation/filesystems/proc.txt#1212
    linux下/proc/stat 計算CPU利用率|Linux內核 - 開源互助社區 http://www.coctec.com/subject/about/185667.html
*/

package procstat

import (
	log "github.com/sirupsen/logrus"
	"github.com/hacking-thursday/sysd/mods"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CPU struct {
	Name      string
	User      int
	Nice      int
	System    int
	Idle      int
	Iowait    int
	Irq       int
	Softirq   int
	Steal     int
	Guest     int
	GuestNice int
}

type ProcStat struct {
	CPUs             []CPU
	CPUTotal         CPU
	Uptime           int
	ProcessesTotal   int
	RunningProcesses int
	BlockedProcesses int
}

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/procstat", get_procstat)
}

func parse_cpu(s []string) CPU {
	var cpu CPU

	cpu.Name = s[0]
	cpu.User, _ = strconv.Atoi(s[1])
	cpu.Nice, _ = strconv.Atoi(s[2])
	cpu.System, _ = strconv.Atoi(s[3])
	cpu.Idle, _ = strconv.Atoi(s[4])
	cpu.Iowait, _ = strconv.Atoi(s[5])
	cpu.Irq, _ = strconv.Atoi(s[6])
	cpu.Softirq, _ = strconv.Atoi(s[7])
	cpu.Steal, _ = strconv.Atoi(s[8])
	cpu.Guest, _ = strconv.Atoi(s[9])
	cpu.GuestNice, _ = strconv.Atoi(s[10])

	return cpu
}

func get_procstat(eng_ifce interface{}, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var output []byte
	var procstat ProcStat
	var cpus []CPU
	var cpu CPU

	path := "/proc/stat"

	ra, _ := ioutil.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(ra)), "\n")

	for _, line := range lines {
		s := strings.Split(line, " ")

		switch s[0] {
		case "btime":
			procstat.Uptime, _ = strconv.Atoi(s[1])
		case "processes":
			procstat.ProcessesTotal, _ = strconv.Atoi(s[1])
		case "procs_running":
			procstat.RunningProcesses, _ = strconv.Atoi(s[1])
		case "procs_blocked":
			procstat.BlockedProcesses, _ = strconv.Atoi(s[1])
		case "cpu":
			procstat.CPUTotal = parse_cpu(s)
		default:
			r, _ := regexp.Compile("cpu.*")
			if r.MatchString(s[0]) {
				cpu = parse_cpu(s)
				cpus = append(cpus, cpu)

				procstat.CPUs = cpus
			}
		}
	}

	if output, err = mods.Marshal(r, procstat); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(output); err != nil {
		mods.HttpError(w, err)
		return
	}

	return
}
