PWD=$(shell pwd)

check_root:
	@if [ $(shell id -u) -ne 0 ]; then \
		echo "This script must be run as root"; \
		exit 1; \
	fi

build:
	./make.sh || true

build2:
	BuildTags="docker" ./make.sh || true

run:
	DEBUG=1 SYSD_UI_DIR="$$(pwd)/mods/ui/files" sysd/sysd

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

setup_init_script: check_root
	@if [ ! -f "/etc/init.d/sysd" ]; then \
		ln -s $(PWD)/debian/sysd.init /etc/init.d/sysd; \
	fi
	update-rc.d sysd defaults

install: check_root setup_init_script
	@if [ ! -f "/usr/bin/sysd" ]; then \
		ln -s $(PWD)/sysd/sysd /usr/bin/sysd; \
	fi
