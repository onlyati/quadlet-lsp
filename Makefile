MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.DEFAULT_GOAL := all
.SUFFIXES:

GO := go
BUILD_ENV := CGO_BUILD=0

# Main targets
.PHONY: all
all: test

.PHONY: test
test:
	@echo "Running all tests..."
	$(GO) test -v -race ./...

.PHONY: test-short
test-short:
	@echo "Running tests (short mode)..."
	$(GO) test -short ./...

.PHONY: test-verbose
test-verbose:
	@echo "Running tests with verbose output..."
	$(GO) test -v ./...

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(GO) test -v -short ./...

.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

.PHONY: lint
lint: fmt vet
	@echo "Linting complete"

.PHONY: build
build: bin/quadlet-lsp
	@echo "Building..."
	$(BUILD_ENV) $(GO) build -mod=vendor -o $<

.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GO) clean -cache -testcache

.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying go.mod..."
	$(GO) mod tidy

.PHONY: mod-verify
mod-verify:
	@echo "Verifying dependencies..."
	$(GO) mod verify

.PHONY: deps
deps: mod-tidy mod-verify
