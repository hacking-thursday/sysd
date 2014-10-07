build:
	GOPATH=$$(pwd) TMPDIR="/tmp" $(MAKE) build_all

build_all:
	# deps for src/api
	go get github.com/docker/libcontainer/user
	go get github.com/docker/docker/engine
	go get github.com/docker/docker/pkg/systemd
	# deps for src/daemon
	go get code.google.com/p/gosqlite/sqlite3
	go get github.com/docker/libtrust
	go get github.com/godbus/dbus
	go get github.com/kr/pty
	go get github.com/syndtr/gocapability/capability
	go get github.com/tchap/go-patricia/patricia
	# check and patch
	test "`md5sum src/github.com/docker/libcontainer/cgroups/systemd/apply_systemd.go | cut -c-7 `" == "4d0aedc" && cp -v misc/apply_systemd.go src/github.com/docker/libcontainer/cgroups/systemd/apply_systemd.go || true
	# start build
	cd src/api/server; go build; go install
	cd src/api; go build; go install
	cd src/builtins; go build; go install
	cd src/daemon; go build; go install
	cd sysd/; go build

run:
	DEBUG=1 sysd/sysd
