.PHONY: test

export GO111MODULE=on

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
VERSION_FLAGS := -X "github.com/deepsourcelabs/revive/cli.version=$(GIT_VERSION)" -X "github.com/deepsourcelabs/revive/cli.date=$(DATE)" -X "github.com/deepsourcelabs/revive/cli.commit=$(GIT_COMMIT)" -X "github.com/deepsourcelabs/revive/cli.builtBy=$(BUILDER)"

install:
	@go mod vendor

build:
	@go build -ldflags='$(VERSION_FLAGS)'

test:
	@go test -v -race ./...

