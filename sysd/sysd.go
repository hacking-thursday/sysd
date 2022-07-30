package main

import (
	flag "flag"
	"sysd/apiserver"
)

func main() {
	flag.Parse()
	apiserver.ListenAndServe(nil)
}
