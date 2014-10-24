FROM golang:1.3.3
WORKDIR /go/src/github.com/hacking-thursday/sysd
ADD . /go/src/github.com/hacking-thursday/sysd
RUN cd sysd && go get -v && go test -v && go build -v
RUN cp sysd/sysd /usr/local/bin
