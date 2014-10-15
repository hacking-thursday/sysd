// +build !docker

package main

import (
	"fmt"
)

func runDaemonByDocker() {
	fmt.Printf("The daemon from docker is not included. Try enable \"docker\" build tag and rebuild if you need.\n")
}
