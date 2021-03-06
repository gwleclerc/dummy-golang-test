.PHONY: build run test lint generate convey

default: run

build:
	go get
	go build

run:
	go get
	go get github.com/skelterjohn/rerun
	rerun github.com/gwleclerc/dummy-golang-test

test:
	go get github.com/smartystreets/goconvey
	go get github.com/AlekSi/gocoverutil
	go get github.com/axw/gocov/...
	go get github.com/AlekSi/gocov-xml
	go get github.com/jstemmer/go-junit-report
	mkdir -p ./dist || true
	gocoverutil -coverprofile=./dist/cover.out test -v -covermode=count github.com/gwleclerc/dummy-golang-test/... | go-junit-report > ./dist/tests.xml
	gocov convert ./dist/cover.out | gocov-xml > ./dist/coverage.xml
	go tool cover -html=./dist/cover.out -o=./dist/cover.html

lint:
	go get github.com/alecthomas/gometalinter
	mkdir -p ./dist || true
	gometalinter --install --force
	gometalinter --fast --tests --vendor --disable=gas --disable=gotype --checkstyle -e mocks ./... > ./dist/checkstyle.xml

convey:
	go get github.com/smartystreets/goconvey
	goconvey ./...

generate:
	go get github.com/vektra/mockery/.../
	mockery -all -dir "./cache" -output "./cache/mocks"

venom:
	go get github.com/ovh/venom/...
	cd ./integration && venom run --details high
