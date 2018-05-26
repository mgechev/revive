deps.devtools:
	@go get github.com/golang/dep/cmd/dep

install: deps.devtools
	@dep ensure -v

build:
	@go build

test.all:
	@go test -v ./test/...

