MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")

GO_PACKAGES = go list ./...
GO_FOLDERS = go list -f '{{ .Dir }}' ./...
GO_FILES = find . -type f -name '*.go'

export GO111MODULE := on
#export GOFLAGS := -mod=vendor
#GOPATH := $(shell go env GOPATH)
GO_VER := $(shell go env GOVERSION)

.PHONY: mod
mod:
	go mod tidy
	go mod verify

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

# https://pkg.go.dev/cmd/go/internal/test
.PHONY: test
test:
	CGO_ENABLED=1 go test -timeout 30s -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-n-read
test-n-read: test
	@go tool cover -func coverage.txt

.PHONY: bench
bench: # runs all benchmarks
	go test -benchmem -run=^Benchmark$$ -mod=readonly -bench=. ./...

.PHONY: ci-bench
ci-bench: # runs all benchmarks
	go test -benchtime=1s -count=7 -benchmem -run=^Benchmark$ -mod=readonly -bench=. ./...
