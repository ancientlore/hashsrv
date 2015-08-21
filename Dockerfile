FROM golang

ADD hashsrv.config /go/etc/hashsrv.config
ADD . /go/src/github.com/ancientlore/hashsrv

RUN go get github.com/ancientlore/flagcfg
RUN go get github.com/facebookgo/flagenv
RUN go get github.com/kardianos/service
RUN go get github.com/golang/snappy
RUN go get golang.org/x/crypto/blowfish
RUN go get golang.org/x/crypto/twofish
RUN go get golang.org/x/crypto/ripemd160

RUN go install github.com/ancientlore/hashsrv/cmd/hashsrv

ENTRYPOINT ["/go/bin/hashsrv", "-run"]

EXPOSE 9009
