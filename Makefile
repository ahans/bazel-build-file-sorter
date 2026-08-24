.PHONY: build test lint vet all

VERSION := $(shell git describe --tags --always --dirty)

all: format build test vet lint

format:
	gofmt -w .

build:
	go build -ldflags "-X main.version=$(VERSION)" ./...

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...
