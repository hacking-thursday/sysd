// +build exclude_graphdriver_aufs

package daemon

import (
	"daemon/graphdriver"
)

func migrateIfAufs(driver graphdriver.Driver, root string) error {
	return nil
}
