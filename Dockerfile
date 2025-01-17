# syntax=docker/dockerfile:1.12
FROM golang:1.23.5 AS builder

SHELL ["/bin/bash", "-c"]

ENV GOMODCACHE=/root/.cache/gomod
ENV CGO_ENABLED=0
ENV REPO=github.com/nabeken/go-api-now

WORKDIR /go/src
COPY . ./

RUN --mount=type=cache,target=/root/.cache \
  set -eo pipefail; \
  go get -d -v ./...; \
  go build -v -o ../bin/go-api-now

FROM gcr.io/distroless/static-debian11:latest

# for debugging
#FROM gcr.io/distroless/static-debian11:debug
#SHELL ["/busybox/sh", "-c"]

COPY --from=builder /go/bin/go-api-now /

EXPOSE 8000

CMD ["/go-api-now"]
