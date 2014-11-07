PWD=$(shell pwd)

build:
	./scripts/make.sh || true

build2:
	BuildTags="docker" ./scripts/make.sh || true

run:
	DEBUG=1 SYSD_UI_DIR="$$(pwd)/mods/ui/files" sysd/sysd

run2:
	DEBUG=1 sysd/sysd --SYSD_BACKEND="docker"

test:
	env

clean:
	rm -rf .gopath || true
	rm -rf .tmp || true

install:
	## @if [ ! -f "/usr/bin/sysd" ]; then \
	## 	ln -s $(PWD)/sysd/sysd /usr/bin/sysd; \
	## fi
	install -D --mode=0644 sysd/sysd $(DESTDIR)/usr/sbin/sysd
