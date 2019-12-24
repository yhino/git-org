export GO111MODULE = on

REVISION := $(shell git rev-parse --short HEAD)

SRCS := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -buildid= -X \"github.com/yhino/git-org/cmd/version.revision=$(REVISION)\" -extldflags \"-static\""


all: deps git-org

git-org: $(SRCS)
	CGO_ENABLED=0 go build -trimpath $(LDFLAGS)

.PHONY: deps

deps:
	go mod vendor

.PHONY: clean

clean:
	rm -f git-org
	rm -rf vendor/*
