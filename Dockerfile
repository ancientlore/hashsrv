FROM golang:1.16 as builder
WORKDIR /go/src/github.com/ancientlore/hashsrv
COPY . .
WORKDIR /go/src/github.com/ancientlore/hashsrv/cmd/hashsrv
RUN go version
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go get .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go install

FROM gcr.io/distroless/static:nonroot
COPY hashsrv.config /go/etc/hashsrv.config
COPY --from=builder /go/bin/hashsrv /go/bin/hashsrv
EXPOSE 9009
ENTRYPOINT ["/go/bin/hashsrv", "-run"]
