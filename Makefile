PACKAGE_NAME := github.com/wheely/tclientpool
IMAGE_NAME := $(shell echo ${PACKAGE_NAME}-$(shell git rev-parse HEAD) | sed -e "s/[\/\.]/-/g")
THRIFT_VERSION := 0.12

build:
	docker build -t ${IMAGE_NAME} .

test:
	docker run --rm -v ${PWD}:/app ${IMAGE_NAME} go test ./...

lint:
	docker run --rm -v ${PWD}:/app ${PACKAGE_NAME} errcheck ${PACKAGE_NAME}
	docker run --rm -v ${PWD}:/app ${PACKAGE_NAME} golint -set_exit_status ${PACKAGE_NAME}

dep:
	docker run --rm -v ${PWD}:/go/src/${PACKAGE_NAME} -w /go/src/${PACKAGE_NAME} instrumentisto/dep:0.4-alpine ensure --vendor-only

gen:
	docker run --rm \
	-v ${PWD}:/data \
	-w /data thrift:${THRIFT_VERSION} \
	thrift -strict --verbose --out . -r --gen "go:package_prefix=${PACKAGE_NAME}/" example.thrift

all: build gen dep test lint
