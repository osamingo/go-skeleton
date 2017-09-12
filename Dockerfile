FROM golang:1.9
WORKDIR /go/src/github.com/osamingo/go-skeleton
ADD . /go/src/github.com/osamingo/go-skeleton
RUN make build

FROM alpine:3.6
RUN apk add --no-cache tzdata ca-certificates
COPY --from=0 /go/src/github.com/osamingo/go-skeleton/build/go-skeleton /go-skeleton
ENTRYPOINT ["/go-skeleton"]

