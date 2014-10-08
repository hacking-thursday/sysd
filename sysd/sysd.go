package main

import (
	flag "github.com/docker/docker/pkg/mflag"

	apiserver "github.com/hacking-thursday/sysd/api/server"
)

func main() {
	flag.Parse()

	apiserver.ListenAndServe()
}
