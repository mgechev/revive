.PHONY: test lint

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
VERSION_FLAGS := -X "github.com/mgechev/revive/cli.version=$(GIT_VERSION)" -X "github.com/mgechev/revive/cli.date=$(DATE)" -X "github.com/mgechev/revive/cli.commit=$(GIT_COMMIT)" -X "github.com/mgechev/revive/cli.builtBy=$(BUILDER)"

all: test lint build

install:
	@go mod vendor

tidy:
	@go mod tidy -diff

build: tidy
	@go build -ldflags='$(VERSION_FLAGS)'

lint:
	revive --config revive.toml ./...
	golangci-lint run

fmt:
	golangci-lint fmt

test:
	@go test -v -race ./...
