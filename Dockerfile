ARG GOLANG_VERSION=1.12
ARG PACKAGE_NAME=github.com/wheely/tclientpool

FROM golang:${GOLANG_VERSION}-alpine

RUN apk add --no-cache git bash musl-dev gcc

RUN go get -u golang.org/x/lint/golint
RUN go get -u github.com/kisielk/errcheck

WORKDIR /go/src/${PACKAGE_NAME}

COPY . .
