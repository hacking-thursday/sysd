package ifconfig

import (
	"net"
	"net/http"

	"github.com/hacking-thursday/sysd/mods"
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
}

func ifconfig(w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
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
