VERSION ?= 1.0.0
VERSION_FILE := $(subst .,-,$(VERSION))
GOOS    ?= $(shell go env GOOS)
GOARCH  ?= $(shell go env GOARCH)
MODULE  := github.com/imflawlezz/cursor-utils/installer
LDFLAGS := -s -w -X $(MODULE)/internal/config.Version=$(VERSION)
CGO     := CGO_ENABLED=0

bin = cursor_utils_installer_$(1)_$(2)_$(VERSION_FILE)

.PHONY: build build-all test tidy clean

build: dist
	$(CGO) go build -ldflags "$(LDFLAGS)" -o dist/$(call bin,$(GOOS),$(GOARCH)) ./cmd/cursor-utils

build-all: dist
	$(CGO) GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(call bin,darwin,arm64) ./cmd/cursor-utils
	$(CGO) GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(call bin,darwin,amd64) ./cmd/cursor-utils
	$(CGO) GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(call bin,linux,arm64) ./cmd/cursor-utils
	$(CGO) GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(call bin,linux,amd64) ./cmd/cursor-utils

dist:
	mkdir -p dist

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf dist
