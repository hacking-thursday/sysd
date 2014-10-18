package ifconfig

import (
    "strings"
    "io/ioutil"
    "regexp"
    "strconv"
//    "fmt"
)

var myregexp = regexp.MustCompile(`(.+?):\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)

func net_io_counters() (iface_counter_map map[string] Counters) {
    ra, _ := ioutil.ReadFile("/proc/net/dev")
    list := strings.Split(strings.TrimSpace(string(ra)),"\n")

    iface_counter_map = make(map[string]Counters)

    for index,value := range list {
        if index > 1 {
            counter := Counters{}

            list_a := myregexp.FindStringSubmatch(strings.TrimSpace(value))

            counter.BytesRecv, _ = strconv.Atoi(list_a[2])
            counter.PacketsRecv, _ = strconv.Atoi(list_a[3])
            counter.Errin, _        = strconv.Atoi(list_a[4])
            counter.Dropin, _       = strconv.Atoi(list_a[5])
            counter.BytesSent, _   = strconv.Atoi(list_a[9])
            counter.PacketsSent, _ = strconv.Atoi(list_a[10])
            counter.Errout, _       = strconv.Atoi(list_a[11])
            counter.Dropout, _      = strconv.Atoi(list_a[12])

            iface_counter_map[list_a[1]] = counter
//            fmt.Println(list_a[1])
        }
    }
//    fmt.Println(iface_counter_map)

    return iface_counter_map
}
