
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
	# check and patch
	@if [ "`md5sum $${GOPATH}/src/github.com/docker/libcontainer/cgroups/systemd/apply_systemd.go | cut -c-7`" = "4d0aedc" ]; then \
		cp -v misc/apply_systemd.go $${GOPATH}/src/github.com/docker/libcontainer/cgroups/systemd/apply_systemd.go; \
	fi
	# start build
	#go build github.com/hacking-thursday/sysd/api/server ; go install github.com/hacking-thursday/sysd/api/server 
	#go build github.com/hacking-thursday/sysd/api        ; go install github.com/hacking-thursday/sysd/api        
	#go build github.com/hacking-thursday/sysd/builtins   ; go install github.com/hacking-thursday/sysd/builtins   
	#go build github.com/hacking-thursday/sysd/daemon     ; go install github.com/hacking-thursday/sysd/daemon     
	#go build github.com/hacking-thursday/sysd/mods       ; go install github.com/hacking-thursday/sysd/mods       
	cd sysd; go build

run:
	DEBUG=1 sysd/sysd
