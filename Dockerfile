FROM golang:1.10.3-alpine3.7 AS build-env

ARG TAG=latest
LABEL tag=$TAG

RUN apk --no-cache add curl git

RUN curl -o dep -sSL https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
    echo '31144e465e52ffbc0035248a10ddea61a09bf28b00784fd3fdd9882c8cbb2315  dep' | sha256sum -c - && \
    mv dep /usr/bin/dep && chmod 755 /usr/bin/dep

RUN mkdir -p /go/src/github.com/bryanl/frontdoor
WORKDIR /go/src/github.com/bryanl/frontdoor

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only
COPY . .
RUN go install github.com/bryanl/frontdoor/cmd/frontdoor

FROM alpine:3.7

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=build-env /go/bin/frontdoor .
ENTRYPOINT ["./frontdoor"]
