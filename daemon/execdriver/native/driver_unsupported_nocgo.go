// +build linux,!cgo

package native

import (
	"fmt"

	"github.com/hacking-thursday/sysd/daemon/execdriver"
)

func NewDriver(root, initPath string) (execdriver.Driver, error) {
	return nil, fmt.Errorf("native driver not supported on non-linux")
}
