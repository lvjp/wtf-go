FROM golang:1.27.2-trixie@sha256:e58d6f83b3416618d8bcac2b3dde1b7f7e3c4a77d25e88637f8bbae81536c48d AS build

WORKDIR /opt/wtf-go/src

RUN go env
RUN [ $(go env GOCACHE) = /root/.cache/go-build ]
RUN [ $(go env GOMODCACHE) = /go/pkg/mod ]

ARG COVER
RUN --mount=type=bind,target=.,readonly \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x && CGO_ENABLED=0 go build ${COVER:+-cover} -o /opt/wtf-go/wtf-go .

FROM debian:13.5-slim@sha256:b6e2a152f22a40ff69d92cb397223c906017e1391a73c952b588e51af8883bf8

COPY --from=build /opt/wtf-go/wtf-go /usr/bin
ENTRYPOINT ["/usr/bin/wtf-go"]
CMD ["serve"]
