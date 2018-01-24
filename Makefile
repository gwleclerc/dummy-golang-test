.PHONY: build run test lint

default: run

build:
	go build

run:
	go get github.com/skelterjohn/rerun
	rerun github.com/gwleclerc/dummy-golang-test

test:
	go test ./...

lint:
	go get github.com/alecthomas/gometalinter
	gometalinter --install --force
	gometalinter --fast --tests --vendor --disable=gas --disable=gotype ./...