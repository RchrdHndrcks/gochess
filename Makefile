SRC := $(wildcard *.go)
GOBIN := $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOBIN)/golangci-lint

.PHONY: all test lint fmt vet tools

all: fmt vet lint test

test:
	@echo "Executing tests..."
	@go test -v ./... | grep "FAIL"

lint: tools
	@echo "Executing linter..."
	@$(GOLANGCI_LINT) run

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Verifying code..."
	@go vet ./...

tools:
	@echo "Installing necessary tools..."
	@if [ ! -f $(GOLANGCI_LINT) ]; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Tools installed successfully."
