FROM golang:1.26.2-trixie@sha256:c0074c718b473f3827043f86532c4c0ff537e3fe7a81b8219b0d1ccfcc2c9a09 AS build

WORKDIR /opt/wtf-go/src

RUN go env
RUN [ $(go env GOCACHE) = /root/.cache/go-build ]
RUN [ $(go env GOMODCACHE) = /go/pkg/mod ]

RUN --mount=type=bind,target=.,readonly \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x && CGO_ENABLED=0 go build -o /opt/wtf-go/wtf-go .

FROM debian:13.4-slim@sha256:109e2c65005bf160609e4ba6acf7783752f8502ad218e298253428690b9eaa4b

COPY --from=build /opt/wtf-go/wtf-go /usr/bin
ENTRYPOINT ["/usr/bin/wtf-go"]
CMD ["serve"]
