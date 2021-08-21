PATH  := $(PATH):$(HOME)/go/bin
SHELL := env PATH=$(PATH) /bin/bash

BIN := "./bin/sysmon"
BUILD_DATE := $(shell date  +%Y-%m-%dT%H:%M)
GIT_HASH := $(shell git log --format="%h" -n 1)
OS := $(shell uname)
ARCH := $(shell uname -m)
LDFLAGS := -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH) -X main.arch=$(ARCH) -X main.osSys=$(OS)


proto:
	protoc -I /usr/local/include -I schema --go-grpc_out=api --go_out=api sysmon.proto

build-server:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/sysmon2