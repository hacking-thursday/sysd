// +build linux

package iptables

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"github.com/docker/docker/pkg/version"
	"github.com/hacking-thursday/sysd/mods"
)

func init() {
	log.Debugf("Initializing module...")
	mods.Register("GET", "/iptables", handler)
}

func handler(eng_ifce interface{}, version version.Version, w http.ResponseWriter, r *http.Request, vars map[string]string) (err error) {
	//result := []map[string]interface{}{}
	result := map[string]interface{}{}

	ip_tables_names_path := "/proc/net/ip_tables_names"
	ip_tables_names := []string{}

	f, _ := os.Open(ip_tables_names_path)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		ip_tables_names = append(ip_tables_names, line)
	}

	for _, table_name := range ip_tables_names {

		res_table := map[string]interface{}{}

		fd, err := socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
		if err != nil {
			return err
		}

		info := ipt_getinfo{}
		copy(info.name[:], table_name)

		info_len := _Socklen(unsafe.Sizeof(info))

		err2 := getsockopt(fd, SOL_IP, IPT_SO_GET_INFO, unsafe.Pointer(&info), &info_len)
		if err2 != nil {
			return err2
		}

		res_table["name"] = string(info.name[:bytes.IndexByte(info.name[:], 0)])
		res_table["valid_hooks"] = info.valid_hooks
		res_table["hook_entry"] = info.hook_entry
		res_table["underflow"] = info.underflow
		res_table["num_entries"] = info.num_entries
		res_table["size"] = info.size

		req_entries_len := _Socklen(int(unsafe.Sizeof(ipt_get_entries{})) + int(info.size))
		req_entries_mem := make([]byte, req_entries_len)

		req_entries := (*ipt_get_entries)(unsafe.Pointer(&req_entries_mem[0]))
		copy(req_entries.name[:], table_name)
		req_entries.size = info.size

		err31 := getsockopt(fd, SOL_IP, IPT_SO_GET_ENTRIES, unsafe.Pointer(req_entries), &req_entries_len)
		if err31 != nil {
			return err31
		}

		entrytable := *(*[10000]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&req_entries.entrytable)) + 4))

		offset_beg := 0
		offset_end := int(req_entries_len) - (int(unsafe.Offsetof(req_entries.entrytable)) + 4)

		curr_chain := string("")
		curr_chain_data := []map[string]interface{}{}
		for pos := offset_beg; pos < offset_end; {
			row_entry := map[string]interface{}{}

			ent := (*ipt_entry)(unsafe.Pointer(&entrytable[pos]))

			if pos+int(ent.next_offset) == offset_end {
				//fmt.Println("NOTE: This entry is last one")
			}

			ip_res := map[string]interface{}{}
			ip_res["src"] = grab_net_ip(&(ent.ip.src))
			ip_res["smsk"] = grab_net_ip(&(ent.ip.smsk))
			ip_res["iniface"] = byte_array_to_string(ent.ip.iniface[:])

			ip_res["dst"] = grab_net_ip(&(ent.ip.dst))
			ip_res["dmsk"] = grab_net_ip(&(ent.ip.dmsk))
			ip_res["outiface"] = byte_array_to_string(ent.ip.outiface[:])
			ip_res["proto"] = ent.ip.proto
			ip_res["flags"] = ent.ip.flags
			ip_res["invflags"] = ent.ip.invflags

			row_entry["ip"] = ip_res
			row_entry["nfcache"] = ent.nfcache
			row_entry["target_offset"] = ent.target_offset
			row_entry["next_offset"] = ent.next_offset
			row_entry["comefrom"] = ent.comefrom

			counter_res := map[string]interface{}{}
			counter_res["pcnt"] = ent.counters.pcnt
			counter_res["bcnt"] = ent.counters.bcnt
			row_entry["counters"] = counter_res

			target := (*xt_entry_target)(unsafe.Pointer(&entrytable[pos+int(ent.target_offset)]))

			if byte_array_to_string(target.u.user.name[:]) == "ERROR" {
				data4 := *(*[32]byte)(unsafe.Pointer(&target.data))
				chain_str := byte_array_to_string(data4[:])
				//fmt.Println("new userdefined chain:", chain_str)

				if curr_chain != chain_str {
					curr_chain_data = []map[string]interface{}{}
				}
				curr_chain = chain_str
				pos += int(ent.next_offset)
				continue
			} else {
				idx := iptcb_ent_is_hook_entry(pos, info)
				if idx != 0 {
					hooknames := [5]string{
						"PREROUTING",
						"INPUT",
						"FORWARD",
						"OUTPUT",
						"POSTROUTING",
					}

					chain_name := hooknames[idx-1]
					//fmt.Println("NOTE: This entry is hook entry, chainname=", chain_name)

					if curr_chain != chain_name {
						curr_chain_data = []map[string]interface{}{}
					}
					curr_chain = chain_name
				} else {

				}
			}

			ptr2 := unsafe.Offsetof(ent.elems)
			for pos2 := 0; int(ptr2)+pos2 < int(ent.target_offset); {
				elem := (*xt_entry_match)(unsafe.Pointer(&entrytable[pos+int(ptr2)+pos2]))

				pos2 += int(elem.u.user.match_size)
			}

			if len(curr_chain) > 0 {
				curr_chain_data = append(curr_chain_data, row_entry)
				res_table[curr_chain] = curr_chain_data
			}

			pos += int(ent.next_offset)
		}

		result[table_name] = res_table
	}

	var out []byte
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

const (
	XT_TABLE_MAXNAMELEN     = 32
	NF_INET_NUMHOOKS        = 5
	IPT_SO_GET_INFO         = 64
	IPT_SO_GET_ENTRIES      = 64 + 1
	SOL_IP                  = 0
	IFNAMSIZ                = 16
	XT_EXTENSION_MAXNAMELEN = 29
	NF_IP_PRE_ROUTING       = 0
	NF_IP_LOCAL_IN          = 1
	NF_IP_FORWARD           = 2
	NF_IP_LOCAL_OUT         = 3
	NF_IP_POST_ROUTING      = 4
)

type _Socklen uint32

type ipt_getinfo struct {
	name        [XT_TABLE_MAXNAMELEN]byte // [0:32]
	valid_hooks uint32                    // [32:36]
	hook_entry  [NF_INET_NUMHOOKS]uint32  // [36:56]
	underflow   [NF_INET_NUMHOOKS]uint32  // [56:76]
	num_entries uint32                    // [76:80]
	size        uint32                    // [80:84]
}

type in_addr struct {
	s_addr uint32
}

type ipt_ip struct {
	src, dst                    in_addr
	smsk, dmsk                  in_addr
	iniface, outiface           [IFNAMSIZ]byte
	iniface_mask, outiface_mask [IFNAMSIZ]uint8
	proto                       uint16
	flags                       uint8
	invflags                    uint8
}

type xt_counters struct {
	pcnt, bcnt uint64
}

type ipt_entry struct {
	ip            ipt_ip      // [0:84]
	nfcache       uint32      // [84:88]
	target_offset uint16      // [88:90]
	next_offset   uint16      // [90:92]
	comefrom      uint32      // [92:96]
	counters      xt_counters // [96:112]
	elems         [4]byte     // pointer
}

type ipt_get_entries struct {
	name       [XT_TABLE_MAXNAMELEN]byte // [0:32]
	size       uint32                    // [32:36]
	entrytable [4]byte                   // pointer to ipt_entry
}

type xt_entry_match struct {
	u struct {
		user struct {
			match_size uint16
			name       [XT_EXTENSION_MAXNAMELEN]byte
			revision   uint8
		}
	}

	data [0]byte
}

type xt_entry_target struct {
	u struct {
		user struct {
			target_size uint16
			name        [XT_EXTENSION_MAXNAMELEN]byte
			revision    uint8
		}
	}

	data [0]byte
}

func getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_GETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen)), 0)
	if e1 != 0 {
		err = e1
	}
	return
}

func socket(domain int, typ int, proto int) (fd int, err error) {
	r0, _, e1 := syscall.RawSyscall(syscall.SYS_SOCKET, uintptr(domain), uintptr(typ), uintptr(proto))
	fd = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func byte_array_to_string(b_ary []byte) (str string) {
	n := bytes.IndexByte(b_ary, 0)
	str = fmt.Sprintf("%s", b_ary[:n])

	return str
}

func grab_net_ip(ip_ptr *in_addr) string {
	ip_bytes := *(*[4]byte)(unsafe.Pointer(ip_ptr))
	res := fmt.Sprintf("%v", net.IPv4(ip_bytes[0], ip_bytes[1], ip_bytes[2], ip_bytes[3]))
	return res
}

func iptcb_ent_is_hook_entry(pos int, info ipt_getinfo) int {

	for i := 0; i < NF_INET_NUMHOOKS; i++ {
		cond1 := (info.valid_hooks&(1<<uint(i)) != 0)
		cond2 := (pos == int(info.hook_entry[i]))

		if cond1 && cond2 {
			return i + 1
		}
	}

	return 0
}
