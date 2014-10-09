// +build !linux

package lxc

import "github.com/hacking-thursday/sysd/daemon/execdriver"

func setHostname(hostname string) error {
	panic("Not supported on darwin")
}

func finalizeNamespace(args *execdriver.InitArgs) error {
	panic("Not supported on darwin")
}
