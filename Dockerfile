FROM golang:1.23 as builder
WORKDIR /go/src/github.com/ancientlore/hashsrv
COPY . .
WORKDIR /go/src/github.com/ancientlore/hashsrv/cmd/hashsrv
RUN go version
RUN CGO_ENABLED=0 GOOS=linux go get .
RUN CGO_ENABLED=0 GOOS=linux go install

FROM ancientlore/goimg:1.23
COPY hashsrv.config /go/etc/hashsrv.config
COPY --from=builder /go/bin/hashsrv /usr/bin/hashsrv
EXPOSE 9009
ENTRYPOINT ["/usr/bin/hashsrv", "-run"]
