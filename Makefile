PWD=$(shell pwd)

build:
	./scripts/make.sh

doc:
	@echo "" > docs/source/modules.rst
	@find mods/ -name "README.rst" -exec echo ".. include:: ../../{}" >> docs/source/modules.rst \;
	cd docs; make html

clean:
	[ -e .gopath ] && rm -rf .gopath || true
	[ -e sysd/sysd ] && rm -rf sysd/sysd || true

install:
	install -D --mode=0644 sysd/sysd $(DESTDIR)/usr/sbin/sysd

.PHONY: clean
