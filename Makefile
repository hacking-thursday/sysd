build:
	./make.sh || true

build2:
	BuildTags="docker" ./make.sh || true

run:
	DEBUG=1 sysd/sysd

run2:
	DEBUG=1 sysd/sysd --SYSD_BACKEND="docker"

test:
	curl "http://127.0.0.1:8080/apilist"
	curl "http://127.0.0.1:8080/ifconfig"
	curl "http://127.0.0.1:8080/info2"
	curl "http://127.0.0.1:8080/loader"
	curl "http://127.0.0.1:8080/memstats"
	curl "http://127.0.0.1:8080/net"
	curl "http://127.0.0.1:8080/osver"
	curl "http://127.0.0.1:8080/ping"
	curl "http://127.0.0.1:8080/sysinfo"
