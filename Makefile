#
# Usage:
# make build           # builds the artifact
# make clean           # removes the artifact and the vendored packages

SHELL := /usr/bin/env bash
BIN_DIR := $(shell pwd)/bin
CMD_DIR := $(shell pwd)/cmd
BIN := dl_cli 
LDFLAGS := ""

.PHONY: clean format test build

default: clean format test build

format:
	gofmt -w -s */*.go

clean:
	rm -rf $(BIN_DIR)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN) main.go

test:
	go test -v ./...

