SHELL := /usr/bin/env bash

.NOTPARALLEL:
.SECONDEXPANSION:
## NOTINTERMEDIATE requires make >=4.4
.NOTINTERMEDIATE:

include $(CURDIR)/scripts/go.mk
include $(CURDIR)/scripts/tools.mk

.DEFAULT_GOAL=default
.PHONY: default
default: checks test

# man git-clean
.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd

.PHONY: env
env:
	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Module \e[0m \e[0;90m<<<\e[0m"
	@echo "$(MODULE)"
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Go env \e[0m \e[0;90m<<<\e[0m"
	go env
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Packages \e[0m \e[0;90m<<<\e[0m"
	$(GO_PACKAGES)
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Folders \e[0m \e[0;90m<<<\e[0m"
	$(GO_FOLDERS)
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Files \e[0m \e[0;90m<<<\e[0m"
	$(GO_FILES)
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Tools \e[0m \e[0;90m<<<\e[0m"
	@echo '$(TOOLS_BIN)'
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Path \e[0m \e[0;90m<<<\e[0m"
	@echo "$${PATH}" | tr ':' '\n'
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Shell \e[0m \e[0;90m<<<\e[0m"
	@echo "SHELL=$${SHELL}"
	@echo "BASH=$${BASH}"
	@echo "BASH_VERSION=$${BASH_VERSION}"
	@echo "BASH_VERSINFO=$${BASH_VERSINFO}"
	@echo ""

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint

.PHONY: ci-format
ci-format: goimports gofumpt
	./scripts/git-check-dirty

.PHONY: ci-mod
ci-mod: mod
	./scripts/git-check-dirty

.PHONY: ci-sh
ci-sh: shfmt shellcheck
	@./scripts/git-check-dirty
