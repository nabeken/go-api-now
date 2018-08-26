FROM golang:1.10.3-alpine
MAINTAINER <nabeken@tknetworks.org>

ENV REPO=github.com/nabeken/go-api-now

RUN apk add --no-cache --update bash git

COPY . src/$REPO
WORKDIR /go/src/$REPO

RUN go get -d -v ./...
RUN go install -v

ENV GIN_MODE=release
EXPOSE 8000
USER nobody
CMD ["go-api-now"]
