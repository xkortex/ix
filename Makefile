VERSION := $(shell git describe --always --dirty --tags)

.PHONY: default get test all

default:
	go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/ix

get:
	go get

fmt:
	go fmt ./...

static: get
	CGO_ENABLED=0 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/ix





