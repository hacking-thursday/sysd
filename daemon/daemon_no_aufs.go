// +build exclude_graphdriver_aufs

package daemon

import (
	"github.com/hacking-thursday/sysd/daemon/graphdriver"
)

func migrateIfAufs(driver graphdriver.Driver, root string) error {
	return nil
}
