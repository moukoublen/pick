SHELL := /bin/bash

.NOTPARALLEL:
.SECONDEXPANSION:
## NOTINTERMEDIATE requires make >=4.4
.NOTINTERMEDIATE:

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
VERSION ?= 0.0.0

PACKAGES = $(GO_EXEC) list -tags=$(TAGS) ./...
FOLDERS = $(GO_EXEC) list -tags=$(TAGS) -f '{{.Dir}}' ./...

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
	$(GO_EXEC) mod tidy -go=1.20
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
	$(PACKAGES)
	@echo ""
	@echo ">>> Folders:"
	$(FOLDERS)
	@echo ""
	@echo ">>> Tools:"
	@echo '$(TOOLSBIN)'
	@echo ">>> Path:"
	@echo $${PATH}


####################################################################################
## <ci & external tools> ###########################################################
####################################################################################
.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="$(TAGS)" -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-n-read
test-n-read: test
	@$(GO_EXEC) tool cover -func coverage.txt

.PHONY: bench
bench:
	CGO_ENABLED=1 $(GO_EXEC) test -benchmem -run=^$$ -mod=readonly -count=1 -v -race -bench=. ./...

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint

.PHONY: vet
vet:
	$(GO_EXEC) vet `$(PACKAGES)`
	@echo ""

TOOLSDIR ?= $(shell pwd)/.ext
TOOLSBIN ?= $(TOOLSDIR)/bin
export TOOLSBIN
export PATH := $(TOOLSBIN):$(PATH)
TOOLSDB ?= $(TOOLSDIR)/.db

uppercase = $(shell echo '$(1)' | tr '[:lower:]' '[:upper:]')

.PHONY: tools
tools: \
	$(TOOLSBIN)/goimports \
	$(TOOLSBIN)/staticcheck \
	$(TOOLSBIN)/golangci-lint \
	$(TOOLSBIN)/gofumpt

$(TOOLSBIN):
	@mkdir -p $(TOOLSBIN)

$(TOOLSDB):
	@mkdir -p $(TOOLSDB)

# In make >= 4.4. .NOTINTERMEDIATE will do the job.
.PRECIOUS: $(TOOLSDB)/%.ver
$(TOOLSDB)/%.ver: | $(TOOLSDB)
	@rm -f $(TOOLSDB)/$(word 1,$(subst ., ,$*)).*
	@touch $(TOOLSDB)/$*.ver

# In make >= 4.4 .NOTINTERMEDIATE will do the job.
.PRECIOUS: $(TOOLSBIN)/%
$(TOOLSBIN)/%: DSC=$*
$(TOOLSBIN)/%: VER=$($(call uppercase,$*)_VER)
$(TOOLSBIN)/%: CMD=$($(call uppercase,$*)_CMD)
$(TOOLSBIN)/%: $(TOOLSDB)/$$(DSC).$$(VER).$(GO_VER).ver
	@echo -e "Installing \e[1;36m$(DSC)\e[0m@\e[1;36m$(VER)\e[0m using \e[1;36m$(GO_VER)\e[0m"
	GOBIN="$(TOOLSBIN)" CGO_ENABLED=0 $(GO_EXEC) install -trimpath -ldflags '-s -w -extldflags "-static"' "$(CMD)@$(VER)"
	@echo ""

## <staticcheck>
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/c
STATICCHECK_CMD=honnef.co/go/tools/cmd/staticcheck
STATICCHECK_VER:=2023.1.6
$(TOOLSDB)/staticcheck.$(STATICCHECK_VER).$(GO_VER).ver:
$(TOOLSBIN)/staticcheck:

.PHONY: staticcheck
staticcheck: $(TOOLSBIN)/staticcheck
	staticcheck -f=text -checks=all,-ST1000 -tests ./...
	@echo ''
## </staticcheck>

## <golangci-lint>
# https://github.com/golangci/golangci-lint/releases
GOLANGCI-LINT_CMD:=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI-LINT_VER:=v1.55.2
$(TOOLSDB)/golangci-lint.$(GOLANGCI-LINT_VER).$(GO_VER).ver:
$(TOOLSBIN)/golangci-lint:

.PHONY: golangci-lint
golangci-lint: $(TOOLSBIN)/golangci-lint
	golangci-lint run
	@echo ''

.PHONY: golangci-lint-github-actions
golangci-lint-github-actions: $(TOOLSBIN)/golangci-lint
	golangci-lint run
	@echo ''
## </golangci-lint>

## <goimports>
# https://pkg.go.dev/golang.org/x/tools?tab=versions
GOIMPORTS_CMD := golang.org/x/tools/cmd/goimports
GOIMPORTS_VER := v0.17.0
$(TOOLSDB)/goimports.$(GOIMPORTS_VER).$(GO_VER).ver:
$(TOOLSBIN)/goimports:

.PHONY: goimports
goimports: $(TOOLSBIN)/goimports
	@echo '$(TOOLSBIN)/goimports -l `$(FOLDERS)`'
	@if [[ -n "$$(goimports -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'goimports errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    make goimports.display'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    make goimports.fix'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: goimports.display
goimports.display: $(TOOLSBIN)/goimports
	goimports -d `$(FOLDERS)`

.PHONY: goimports.fix
goimports.fix: $(TOOLSBIN)/goimports
	goimports -w `$(FOLDERS)`
## </goimports>

## <gofumpt>
# https://github.com/mvdan/gofumpt/releases
GOFUMPT_CMD:=mvdan.cc/gofumpt
GOFUMPT_VER:=v0.5.0
$(TOOLSDB)/gofumpt.$(GOFUMPT_VER).$(GO_VER).ver:
$(TOOLSBIN)/gofumpt:

.PHONY: gofumpt
gofumpt: $(TOOLSBIN)/gofumpt
	@echo '$(TOOLSBIN)/gofumpt -l `$(FOLDERS)`'
	@if [[ -n "$$(gofumpt -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'gofumpt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    make gofumpt.display'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    make gofumpt.fix'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofumpt.display
gofumpt.display:
	gofumpt -d `$(FOLDERS)`

.PHONY: gofumpt.fix
gofumpt.fix:
	gofumpt -w `$(FOLDERS)`
## </gofumpt>

## <gofmt>
.PHONY: gofmt
gofmt:
	@echo 'gofmt -l `$(FOLDERS)`'
	@if [[ -n "$$(gofmt -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'gofmt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    make gofmt.display'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    make gofmt.fix'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofmt.display
gofmt.display:
	gofmt -d `$(FOLDERS)`

.PHONY: gofmt.fix
gofmt.fix:
	gofmt -w `$(FOLDERS)`
## </gofmt>
####################################################################################
## </ci & external tools> ##########################################################
####################################################################################


# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables
# https://www.gnu.org/software/make/manual/make.html#Prerequisite-Types

