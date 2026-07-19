.PHONY: vet
vet:
	go vet `$(GO_PACKAGES)`
	@echo ""

.PHONY: golangci-lint
golangci-lint:
	golangci-lint -v run
	@echo ''

.PHONY: golangci-lint-fmt
golangci-lint-fmt:
	golangci-lint fmt
	@echo ''

.PHONY: gofmt
gofmt:
	gofmt -w `$(GO_FILES)`

.PHONY: gofmt.display
gofmt.display:
	gofmt -d `$(GO_FILES)`

.PHONY: shfmt
shfmt:
	./scripts/foreach-script shfmt \
		--simplify \
		--language-dialect auto \
		--case-indent \
		--indent 2 \
		--write

.PHONY: shellcheck
shellcheck:
	./scripts/foreach-script shellcheck \
		--norc \
		--external-sources \
		--format=tty \
		--enable=require-variable-braces,add-default-case
