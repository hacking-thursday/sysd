build: 
	GOPATH=$$(pwd) TMPDIR="/tmp" $(MAKE) build_all

build_all:
	go get code.google.com/p/go.net/websocket
	go get github.com/gorilla/mux
	go get github.com/docker/libcontainer/user                                                                      
	go get github.com/docker/docker/api
	go get github.com/docker/docker/engine
	go get github.com/docker/docker/pkg/listenbuffer
	go get github.com/docker/docker/pkg/log
	go get github.com/docker/docker/pkg/parsers
	go get github.com/docker/docker/pkg/stdcopy
	go get github.com/docker/docker/pkg/systemd
	go get github.com/docker/docker/pkg/version
	go get github.com/docker/docker/registry
	go get github.com/docker/docker/utils
	go get code.google.com/p/gosqlite/sqlite3
	go get github.com/docker/libtrust
	go get github.com/docker/libtrust/trustgraph
	go get github.com/godbus/dbus
	go get github.com/kr/pty
	go get github.com/syndtr/gocapability/capability
	go get github.com/tchap/go-patricia/patricia
	#go get github.com/docker/libcontainer/cgroups/systemd
	cd src/api/server; go build; go install
	cd src/api; go build; go install
	cd src/builtins; go build; go install
	cd sysd/; go build

run:
	sysd/sysd
