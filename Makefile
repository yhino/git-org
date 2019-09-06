export GO111MODULE = on

REVISION := $(shell git rev-parse --short HEAD)

SRCS := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"github.com/yhino/git-org/cmd/version.revision=$(REVISION)\" -extldflags \"-static\""


all: deps git-org

git-org: $(SRCS)
	go build $(LDFLAGS)

.PHONY: deps

deps:
	go mod vendor

.PHONY: clean

clean:
	rm -f git-org
	rm -rf vendor/*
