MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")

GO_PACKAGES = go list -tags='$(TAGS)' ./...
GO_FOLDERS = go list -tags='$(TAGS)' -f '{{ .Dir }}' ./...
GO_FILES = find . -type f -name '*.go' -not -path './vendor/*'

export GO111MODULE := on
#export GOFLAGS := -mod=vendor
#GOPATH := $(shell go env GOPATH)
GO_VER := $(shell go env GOVERSION)

.PHONY: mod
mod:
	go mod tidy -go=1.22
	go mod verify

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

.PHONY: test
test:
	CGO_ENABLED=1 go test -timeout 60s -race -tags='$(TAGS)' -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-n-read
test-n-read: test
	@go tool cover -func coverage.txt

.PHONY: bench
bench: # runs all benchmarks
	CGO_ENABLED=1 go test -benchmem -run=^Benchmark$$ -mod=readonly -count=1 -v -race -bench=. ./...
