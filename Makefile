GO ?= go
GOCACHE ?= $(PWD)/.gocache
GOMODCACHE ?= $(PWD)/.gocache/mod
CACHE_ENV = GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)

.PHONY: all fmt test build plugin integ integ-message clean

all: test

fmt:
	$(GO)fmt ./...

test:
	$(CACHE_ENV) $(GO) test ./...

build:
	$(CACHE_ENV) $(GO) build ./...

plugin:
	$(CACHE_ENV) $(GO) build -o bin/messagingplugin ./cmd/messagingplugin

integ-message:
	$(CACHE_ENV) $(GO) run ./integ/messaging.go

integ: integ-message

clean:
	rm -rf $(GOCACHE) bin
