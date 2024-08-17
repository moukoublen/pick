# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables
# https://www.gnu.org/software/make/manual/make.html#Prerequisite-Types

TOOLS_DIR ?= $(shell pwd)/.tools
TOOLS_DB ?= $(TOOLS_DIR)/.db
TOOLS_BIN ?= $(TOOLS_DIR)/bin
export TOOLS_BIN
export PATH := $(TOOLS_BIN):$(PATH)

.PHONY: tools
tools: \
	$(TOOLS_BIN)/goimports \
	$(TOOLS_BIN)/staticcheck \
	$(TOOLS_BIN)/golangci-lint \
	$(TOOLS_BIN)/gofumpt \
	$(TOOLS_BIN)/gojq

.PHONY: clean-tools
clean-tools:
	rm -rf $(TOOLS_DIR)

$(TOOLS_BIN):
	@mkdir -p $(TOOLS_BIN)

$(TOOLS_DB):
	@mkdir -p $(TOOLS_DB)

# In make >= 4.4. .NOTINTERMEDIATE will do the job.
.PRECIOUS: $(TOOLS_DB)/%.ver
$(TOOLS_DB)/%.ver: | $(TOOLS_DB)
	@rm -f $(TOOLS_DB)/$(word 1,$(subst ., ,$*)).*
	@touch $(TOOLS_DB)/$*.ver

define go_install
	@echo -e "Installing \e[1;36m$(1)\e[0m@\e[1;36m$(3)\e[0m using \e[1;36m$(GO_VER)\e[0m"
	GOBIN="$(TOOLS_BIN)" CGO_ENABLED=0 $(GO_EXEC) install -trimpath -ldflags '-s -w -extldflags "-static"' "$(2)@$(3)"
	@echo ""
endef

.PHONY: vet
vet:
	$(GO_EXEC) vet `$(GO_PACKAGES)`
	@echo ""

## <staticcheck>
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/c
STATICCHECK_CMD:=honnef.co/go/tools/cmd/staticcheck
STATICCHECK_VER:=2024.1.1
$(TOOLS_BIN)/staticcheck: $(TOOLS_DB)/staticcheck.$(STATICCHECK_VER).$(GO_VER).ver
	$(call go_install,staticcheck,$(STATICCHECK_CMD),$(STATICCHECK_VER))

.PHONY: staticcheck
staticcheck: $(TOOLS_BIN)/staticcheck
	$(TOOLS_BIN)/staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''
## </staticcheck>

## <golangci-lint>
# https://github.com/golangci/golangci-lint/releases
GOLANGCI-LINT_CMD:=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI-LINT_VER:=v1.60.1
$(TOOLS_BIN)/golangci-lint: $(TOOLS_DB)/golangci-lint.$(GOLANGCI-LINT_VER).$(GO_VER).ver
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_BIN) $(GOLANGCI-LINT_VER)

.PHONY: golangci-lint
golangci-lint: $(TOOLS_BIN)/golangci-lint
	$(TOOLS_BIN)/golangci-lint run
	@echo ''

.PHONY: golangci-lint-github-actions
golangci-lint-github-actions: $(TOOLS_BIN)/golangci-lint
	golangci-lint run --out-format github-actions
	@echo ''
## </golangci-lint>

## <goimports>
# https://pkg.go.dev/golang.org/x/tools?tab=versions
GOIMPORTS_CMD := golang.org/x/tools/cmd/goimports
GOIMPORTS_VER := v0.24.0
$(TOOLS_BIN)/goimports: $(TOOLS_DB)/goimports.$(GOIMPORTS_VER).$(GO_VER).ver
	$(call go_install,goimports,$(GOIMPORTS_CMD),$(GOIMPORTS_VER))

.PHONY: goimports
goimports: $(TOOLS_BIN)/goimports
	@echo '$(TOOLS_BIN)/goimports -l `$(GO_FILES)`'
	@if [[ -n "$$($(TOOLS_BIN)/goimports -l `$(GO_FILES)` | tee /dev/stderr)" ]]; then \
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
goimports.display: $(TOOLS_BIN)/goimports
	$(TOOLS_BIN)/goimports -d `$(GO_FOLDERS)`

.PHONY: goimports.fix
goimports.fix: $(TOOLS_BIN)/goimports
	$(TOOLS_BIN)/goimports -w `$(GO_FOLDERS)`
## </goimports>

## <gofumpt>
# https://github.com/mvdan/gofumpt/releases
GOFUMPT_CMD:=mvdan.cc/gofumpt
GOFUMPT_VER:=v0.7.0
$(TOOLS_BIN)/gofumpt: $(TOOLS_DB)/gofumpt.$(GOFUMPT_VER).$(GO_VER).ver
	$(call go_install,gofumpt,$(GOFUMPT_CMD),$(GOFUMPT_VER))

.PHONY: gofumpt
gofumpt: $(TOOLS_BIN)/gofumpt
	@echo '$(TOOLS_BIN)/gofumpt -l `$(GO_FOLDERS)`'
	@if [[ -n "$$($(TOOLS_BIN)/gofumpt -l `$(GO_FOLDERS)` | tee /dev/stderr)" ]]; then \
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
	$(TOOLS_BIN)/gofumpt -d `$(GO_FOLDERS)`

.PHONY: gofumpt.fix
gofumpt.fix:
	$(TOOLS_BIN)/gofumpt -w `$(GO_FOLDERS)`
## </gofumpt>

## <gofmt>
.PHONY: gofmt
gofmt:
	@echo 'gofmt -l `$(GO_FOLDERS)`'
	@if [[ -n "$$(gofmt -l `$(GO_FOLDERS)` | tee /dev/stderr)" ]]; then \
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
	gofmt -d `$(GO_FOLDERS)`

.PHONY: gofmt.fix
gofmt.fix:
	gofmt -w `$(GO_FOLDERS)`
## </gofmt>

## <gojq>
# https://github.com/itchyny/gojq/releases
GOJQ_CMD := github.com/itchyny/gojq/cmd/gojq
GOJQ_VER := v0.12.16
$(TOOLS_BIN)/gojq: $(TOOLS_DB)/gojq.$(GOJQ_VER).$(GO_VER).ver
	$(call go_install,gojq,$(GOJQ_CMD),$(GOJQ_VER))

.PHONY: gojq
gojq: $(TOOLS_BIN)/gojq
## </gojq>
