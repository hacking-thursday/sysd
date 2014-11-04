PWD=$(shell pwd)

build:
	./make.sh || true

build2:
	BuildTags="docker" ./make.sh || true

run:
	DEBUG=1 SYSD_UI_DIR="$$(pwd)/mods/ui/files" sysd/sysd

run2:
	DEBUG=1 sysd/sysd --SYSD_BACKEND="docker"

test:
	env

setup_init_script: check_root
	@if [ ! -f "/etc/init.d/sysd" ]; then \
		ln -s $(PWD)/debian/sysd.init /etc/init.d/sysd; \
	fi
	update-rc.d sysd defaults

check_root:
	@if [ $(shell id -u) -ne 0 ]; then \
		echo "This script must be run as root"; \
		exit 1; \
	fi

install: check_root setup_init_script
	## @if [ ! -f "/usr/bin/sysd" ]; then \
	## 	ln -s $(PWD)/sysd/sysd /usr/bin/sysd; \
	## fi
	install -D --mode=0644 sysd/sysd $(DESTDIR)/usr/sbin/sysd
