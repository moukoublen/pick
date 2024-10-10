SHELL := /bin/bash

.NOTPARALLEL:
.SECONDEXPANSION:
## NOTINTERMEDIATE requires make >=4.4
.NOTINTERMEDIATE:

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
VERSION ?= 0.0.0

GO_PACKAGES = $(GO_EXEC) list -tags='$(TAGS)' ./...
GO_FOLDERS = $(GO_EXEC) list -tags='$(TAGS)' -f '{{ .Dir }}' ./...
GO_FILES = find . -type f -name '*.go' -not -path './vendor/*'

export GO111MODULE := on
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)
GO_VER := $(shell go env GOVERSION)

GO_EXEC ?= go
export GO_EXEC
DOCKER_EXEC ?= docker
export DOCKER_EXEC

.DEFAULT_GOAL=default
.PHONY: default
default: checks test

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -go=1.22
	$(GO_EXEC) mod verify

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

# man git-clean
.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd

.PHONY: env
env:
	@echo "Module: $(MODULE)"
	$(GO_EXEC) env
	@echo ""
	@echo ">>> Packages:"
	$(GO_PACKAGES)
	@echo ""
	@echo ">>> Folders:"
	$(GO_FOLDERS)
	@echo ""
	@echo ">>> Tools:"
	@echo '$(TOOLS_BIN)'
	@echo ""
	@echo ">>> Path:"
	@echo "$${PATH}" | tr ':' '\n'

.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags='$(TAGS)' -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-n-read
test-n-read: test
	@$(GO_EXEC) tool cover -func coverage.txt

.PHONY: bench
bench:
	CGO_ENABLED=1 $(GO_EXEC) test -benchmem -run=^$$ -mod=readonly -count=1 -v -race -bench=. ./...

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint

include $(CURDIR)/scripts/tools.mk

.PHONY: ci-format
ci-format: goimports gofumpt
	./scripts/git-check-dirty

.PHONY: ci-mod
ci-mod: mod
	./scripts/git-check-dirty

.PHONY: ci-sh
ci-sh: shfmt
	@./scripts/sh-checks
	@./scripts/git-check-dirty
