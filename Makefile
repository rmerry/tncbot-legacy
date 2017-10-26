go_packages := $(shell go list ./... | grep -v /vendor/)

.PHONY: default lint test build

default: run

run:
	go run -race tncbot.go -logtostderr=true -config config.json

lint:
	@golint -set_exit_status $(go_packages)

test: lint
	@go test -race $(go_packages)

build: test
	GOARCH=amd64 GOOS=linux go build
