build:
	./make.sh || true

run:
	DEBUG=1 sysd/sysd

test:
	curl "http://127.0.0.1:8080/ping"
	curl "http://127.0.0.1:8080/info2"
