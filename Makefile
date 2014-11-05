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

install:
	## @if [ ! -f "/usr/bin/sysd" ]; then \
	## 	ln -s $(PWD)/sysd/sysd /usr/bin/sysd; \
	## fi
	install -D --mode=0644 sysd/sysd $(DESTDIR)/usr/sbin/sysd
