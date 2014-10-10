package main

import (
	flag "github.com/docker/docker/pkg/mflag"
)

var (
	flBackend = flag.String(
		[]string{"-SYSD_BACKEND"},
		"",
		"Sysd Backend ( docker | minimal ) ",
	)
)
