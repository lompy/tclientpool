PACKAGE_NAME := github.com/wheely/tclientpool

DOCKER_RUN := docker run \
	--rm \
	--user=$(shell id -u):$(shell id -g) \
	--volume ${PWD}:/go/src/${PACKAGE_NAME} \
	--workdir /go/src/${PACKAGE_NAME} \
	--env GOFLAGS=-mod=vendor \
	--env GOCACHE=/tmp/gocache/ \
	--env GO111MODULE=on \
	${PACKAGE_NAME}

build:
	docker build -t ${PACKAGE_NAME} --build-arg PACKAGE_NAME=${PACKAGE_NAME} .

test:
	${DOCKER_RUN} test ./...

lint:
	${DOCKER_RUN} errcheck ${PACKAGE_NAME}
	${DOCKER_RUN} golint -set_exit_status ${PACKAGE_NAME}

dep:
	$(DOCKER_RUN) go mod vendor

gen:
	docker run --rm \
	-v ${PWD}:/data \
	--user=$(shell id -u) \
	-w /data thrift:0.11 \
	thrift -strict --verbose --out . -r --gen "go:package_prefix=${PACKAGE_NAME}/,thrift_import=github.com/apache/thrift/lib/go/thrift" example.thrift

all: build gen dep test lint
