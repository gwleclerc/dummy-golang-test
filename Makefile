.PHONY: build run

default: run

build:
	go build

run:
	go get github.com/skelterjohn/rerun
	rerun github.com/gwleclerc/dummy-golang-test
