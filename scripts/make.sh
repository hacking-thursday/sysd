#!/usr/bin/env bash

set -e

(
	cd $(dirname $0)/../sysd || exit

	export GOPATH="$(pwd)/.gopath:$(pwd)/vendor"
	go get
	CGO_ENABLED=0 go build -v -ldflags="-extldflags=-static"
)
