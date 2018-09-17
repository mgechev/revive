.PHONY: test

install:
	@go install

build:
	@go build

test:
	@go test -v ./test/...

