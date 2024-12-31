# Makefile for zaplint

all: build

build:
	go build -o bin/zaplint cmd/main.go

test-deps:
	cd testdata/src && go mod vendor && cd ../..

test: test-deps
	go test ./...

test-coverage: test-deps
	go test ./... -coverprofile=coverage.out

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

.PHONY: all build tidy lint test-deps test test-coverage
