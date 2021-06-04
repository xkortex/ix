VERSION := $(shell git describe --always --dirty --tags)

.PHONY: default get test fmt all

default:
	go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/ix

get:
	go get

fmt:
	go fmt ./...

test:
	go test ./...

static: get
	CGO_ENABLED=0 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/ix





