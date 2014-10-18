package ifconfig

import (
	"net"
	"net/http"

	"github.com/docker/docker/pkg/version"

	"github.com/hacking-thursday/sysd/mods"
//        "fmt"
)

func init() {
	mods.Register("GET", "/ifconfig", ifconfig)
}

type iface_t struct {
	Index  int
	MTU    int
	Name   string
	IP     []string
	HwAddr string
	Flag   struct {
		Up           bool
		Broadcast    bool
		Loopback     bool
		PointToPoint bool
		Multicast    bool
	}
        Counters Counters
}

type Counters struct {
        BytesRecv   int
        PacketsRecv int
        Errin        int
        Dropin       int
        BytesSent   int
        PacketsSent int
        Errout       int
        Dropout      int
}

func ifconfig(engine interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	var (
		out       []byte
		outIfaces []iface_t
		outIface  iface_t
		ifaces    []net.Interface
		iface     net.Interface
		addrs     []net.Addr
		addr      net.Addr
		ip        string
	)

	if ifaces, err = net.Interfaces(); err != nil {
		mods.HttpError(w, err)
		return
	}

        iface_counter_map := net_io_counters()

	for _, iface = range ifaces {
		outIface = iface_t{
			Index:  iface.Index,
			MTU:    iface.MTU,
			Name:   iface.Name,
			HwAddr: iface.HardwareAddr.String(),
		}

		outIface.Flag.Up = iface.Flags&net.FlagUp != 0
		outIface.Flag.Broadcast = iface.Flags&net.FlagBroadcast != 0
		outIface.Flag.Loopback = iface.Flags&net.FlagLoopback != 0
		outIface.Flag.PointToPoint = iface.Flags&net.FlagPointToPoint != 0
		outIface.Flag.Multicast = iface.Flags&net.FlagMulticast != 0


		addrs, err = iface.Addrs()
		if err != nil {
			mods.HttpError(w, err)
			return
		}

		for _, addr = range addrs {
			ip = addr.String()
			outIface.IP = append(outIface.IP, ip)
		}

                outIface.Counters = iface_counter_map[iface.Name]
//                fmt.Println(iface_counter_map[iface.Name])

		outIfaces = append(outIfaces, outIface)
	}

	if out, err = mods.Marshal(r, outIfaces); err != nil {
		mods.HttpError(w, err)
		return
	}

	if _, err = w.Write(out); err != nil {
		mods.HttpError(w, err)
		return
	}
	return
}
