#!/usr/bin/env bash

set -e

(
	cd $(dirname $0)/../sysd || exit

	export GOPATH="$(pwd)/.gopath:$(pwd)/vendor"
	go get
	go build -v
)
