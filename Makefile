PWD=$(shell pwd)

build:
	./scripts/make.sh || true

run:
	DEBUG=1 SYSD_UI_DIR="$$(pwd)/mods/ui/files" sysd/sysd

test:
	env

doc:
	@echo "" > docs/source/modules.rst
	@find mods/ -name "README.rst" -exec echo ".. include:: ../../{}" >> docs/source/modules.rst \;
	cd docs; make html

clean:
	rm -rf .gopath || true
	rm -rf .tmp || true
	rm -rf sysd/sysd || true

install:
	install -D --mode=0644 sysd/sysd $(DESTDIR)/usr/sbin/sysd
	install -D --mode=0644 sysd-cli/sysd-cli $(DESTDIR)/usr/bin/sysd-cli

.PHONY: clean
