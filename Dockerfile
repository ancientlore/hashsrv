FROM golang

ADD hashsrv.config /go/etc/hashsrv.config
ADD . /go/src/github.com/ancientlore/hashsrv

WORKDIR /go/src/github.com/ancientlore/hashsrv/cmd/hashsrv

RUN go get
RUN go install

WORKDIR /go

ENTRYPOINT ["/go/bin/hashsrv", "-run"]

EXPOSE 9009
