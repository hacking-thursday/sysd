
build:
	GOPATH="$$(pwd)/.gopath" TMPDIR="/tmp" $(MAKE) build_all

build_all:
	# deps for src/api
	go get github.com/docker/libcontainer/user
	go get github.com/docker/docker/engine
	go get github.com/docker/docker/pkg/systemd
	go get code.google.com/p/go.net/websocket
	go get github.com/gorilla/mux
	# deps for src/daemon
	go get code.google.com/p/gosqlite/sqlite3
	go get github.com/docker/libtrust
	go get github.com/godbus/dbus
	go get github.com/kr/pty
	go get github.com/syndtr/gocapability/capability
	go get github.com/tchap/go-patricia/patricia
	# linking
	mkdir -p $${GOPATH}/src/github.com/hacking-thursday/; cd $${GOPATH}/src/github.com/hacking-thursday/ ; ( ln -s ../../../../ sysd 2>/dev/null || true )
	# check and patch
	@if [ "`md5sum $${GOPATH}/src/github.com/docker/libcontainer/cgroups/systemd/apply_systemd.go | cut -c-7`" = "4d0aedc" ]; then \
		cp -v misc/apply_systemd.go $${GOPATH}/src/github.com/docker/libcontainer/cgroups/systemd/apply_systemd.go; \
	fi
	# start build
	cd sysd; 

run:
	DEBUG=1 sysd/sysd
