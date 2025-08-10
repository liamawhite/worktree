# Copyright 2025 Worktree Authors
# Licensed under the Apache License, Version 2.0

export NIX_CONFIG := warn-dirty = false

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -w -s -X main.version=$(VERSION) \
           -X main.commit=$(COMMIT) \
           -X main.date=$(DATE)

.PHONY: build check clean format lint test-unit dirty

check: format lint test-unit dirty

format:
	licenser apply -r "Liam White"
	gofmt -w .

lint:
	golangci-lint run

test:
	go test -race -v ./...

build:
	@mkdir -p bin
	@go build -ldflags "$(LDFLAGS)" -o bin/wt .

clean:
	@rm -rf bin/

dirty:
	@git diff --exit-code > /dev/null || (echo "Working directory is dirty" && exit 1)

