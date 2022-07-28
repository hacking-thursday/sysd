FROM golang:1.3.3

ADD . /go/src/github.com/hacking-thursday/sysd

WORKDIR /go/src/github.com/hacking-thursday/sysd/sysd

RUN go get -v && go test -v && go build -v
