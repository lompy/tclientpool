PACKAGE_NAME := tclientpool

build:
	docker build -t ${PACKAGE_NAME} .

test:
	docker run --rm -v ${PWD}:/go/src/${PACKAGE_NAME} ${PACKAGE_NAME} go test ./...

lint:
	docker run --rm -v ${PWD}:/go/src/${PACKAGE_NAME} ${PACKAGE_NAME} errcheck ${PACKAGE_NAME}
	docker run --rm -v ${PWD}:/go/src/${PACKAGE_NAME} ${PACKAGE_NAME} golint -set_exit_status ${PACKAGE_NAME}

dep:
	docker run --rm -v ${PWD}:/go/src/${PACKAGE_NAME} -w /go/src/${PACKAGE_NAME} instrumentisto/dep:0.4-alpine ensure --vendor-only

gen:
	docker run --rm \
	-v ${PWD}:/data \
	-w /data thrift:0.11 \
	thrift -strict --verbose --out . -r --gen "go:package_prefix=${PACKAGE_NAME}/,thrift_import=github.com/apache/thrift/lib/go/thrift" example.thrift

all: build gen dep test lint
