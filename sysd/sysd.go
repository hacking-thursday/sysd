package main

import (
	flag "flag"

	apiserver "github.com/hacking-thursday/sysd/api/server2"
)

func main() {
	flag.Parse()
	apiserver.ListenAndServe(nil)
}
