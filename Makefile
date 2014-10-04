dummy: build_all

all: 
	GOPATH=$$(pwd) $(MAKE) build_all

build_all:
	go get code.google.com/p/go.net/websocket
	go get github.com/gorilla/mux
	cd src/api/server; go build

