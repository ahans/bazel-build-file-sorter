.PHONY: build test lint vet all

all: format build test vet lint

format:
	gofmt -w .

build:
	go build ./...

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...
