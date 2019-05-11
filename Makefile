REVISION := $(shell git rev-parse --short HEAD)

SRCS := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"github.com/yhino/git-org/cmd/version.revision=$(REVISION)\" -extldflags \"-static\""


all: deps git-org

git-org: $(SRCS)
	go build $(LDFLAGS)

.PHONY: deps

deps:
ifeq ($(shell which dep),)
	go get -u github.com/golang/dep/cmd/dep
	hash -r
endif
	dep ensure

.PHONY: clean

clean:
	rm -f git-org
	rm -rf vendor/*
