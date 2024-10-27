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
	$(TOOLS_BIN)/gojq \
	$(TOOLS_BIN)/shfmt

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
GOLANGCI-LINT_VER:=v1.61.0
$(TOOLS_BIN)/golangci-lint: $(TOOLS_DB)/golangci-lint.$(GOLANGCI-LINT_VER).$(GO_VER).ver
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_BIN) $(GOLANGCI-LINT_VER)

.PHONY: golangci-lint
golangci-lint: $(TOOLS_BIN)/golangci-lint
	golangci-lint run --out-format colored-line-number
	@echo ''
## </golangci-lint>

## <goimports>
# https://pkg.go.dev/golang.org/x/tools?tab=versions
GOIMPORTS_CMD := golang.org/x/tools/cmd/goimports
GOIMPORTS_VER := v0.26.0
$(TOOLS_BIN)/goimports: $(TOOLS_DB)/goimports.$(GOIMPORTS_VER).$(GO_VER).ver
	$(call go_install,goimports,$(GOIMPORTS_CMD),$(GOIMPORTS_VER))

.PHONY: goimports
goimports: $(TOOLS_BIN)/goimports
	goimports -w `$(GO_FILES)`

.PHONY: goimports.display
goimports.display: $(TOOLS_BIN)/goimports
	goimports -d `$(GO_FILES)`
## </goimports>

## <gofumpt>
# https://github.com/mvdan/gofumpt/releases
GOFUMPT_CMD:=mvdan.cc/gofumpt
GOFUMPT_VER:=v0.7.0
$(TOOLS_BIN)/gofumpt: $(TOOLS_DB)/gofumpt.$(GOFUMPT_VER).$(GO_VER).ver
	$(call go_install,gofumpt,$(GOFUMPT_CMD),$(GOFUMPT_VER))

.PHONY: gofumpt
gofumpt: $(TOOLS_BIN)/gofumpt
	gofumpt -w `$(GO_FILES)`

.PHONY: gofumpt.display
gofumpt.display:
	gofumpt -d `$(GO_FILES)`
## </gofumpt>

## <gofmt>
.PHONY: gofmt
gofmt:
	gofmt -s -w `$(GO_FILES)`

.PHONY: gofmt.display
gofmt.display:
	gofmt -d `$(GO_FILES)`
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

## <shfmt>
# https://github.com/mvdan/sh/releases
SHFMT_CMD := mvdan.cc/sh/v3/cmd/shfmt
SHFMT_VER := v3.10.0
$(TOOLS_BIN)/shfmt: $(TOOLS_DB)/shfmt.$(SHFMT_VER).$(GO_VER).ver
	$(call go_install,shfmt,$(SHFMT_CMD),$(SHFMT_VER))

.PHONY: shfmt
shfmt: $(TOOLS_BIN)/shfmt
	@VERBOSE=1 ./scripts/foreach-script shfmt \
	--simplify \
	--language-dialect auto \
	--case-indent \
	--indent 2 \
	--write
## <shfmt>

## <shellcheck>
# https://github.com/koalaman/shellcheck/releases
SHELLCHECK_VER := v0.10.0
$(TOOLS_BIN)/shellcheck: $(TOOLS_DB)/shellcheck.$(SHELLCHECK_VER).ver | $(TOOLS_BIN)
	#./scripts/install-shellcheck --version $(SHELLCHECK_VER) --destination $(TOOLS_DIR)
	$(GO_EXEC) run -tags=install_shellcheck scripts/install-shellcheck.go

.PHONY: shellcheck
shellcheck: $(TOOLS_BIN)/shellcheck
	@VERBOSE=1 ./scripts/foreach-script $(TOOLS_BIN)/shellcheck --external-sources --format=tty --severity=info
## </shellcheck>

