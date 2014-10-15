package main

import (
	flag "github.com/docker/docker/pkg/mflag"

	apiserver "github.com/hacking-thursday/sysd/api/server2"
)

func runDaemonByKdTsai() {
	apiserver.ListenAndServe()
}

func main() {
	flag.Parse()

	if *flBackend == "minimal" {
		runDaemonByKdTsai()
	} else if *flBackend == "docker" {
		runDaemonByDocker()
	} else {
		runDaemonByKdTsai()
	}
}
