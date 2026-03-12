FROM golang:1.26-alpine3.23@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS build

WORKDIR /opt/wtf-go/src

RUN go env
RUN [ $(go env GOCACHE) = /root/.cache/go-build ]
RUN [ $(go env GOMODCACHE) = /go/pkg/mod ]

RUN --mount=type=bind,target=.,readonly \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x && CGO_ENABLED=0 go build -o /opt/wtf-go/wtf-go .

FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

COPY --from=build /opt/wtf-go/wtf-go /usr/bin
ENTRYPOINT ["/usr/bin/wtf-go"]
CMD ["serve"]
