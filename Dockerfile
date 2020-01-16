FROM golang:1.13-alpine

RUN apk add --no-cache git bash musl-dev gcc

RUN go get -u golang.org/x/lint/golint
RUN go get -u github.com/kisielk/errcheck

WORKDIR /app

COPY . .
