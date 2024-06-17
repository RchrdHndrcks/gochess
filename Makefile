SRC := $(wildcard *.go)
GOBIN := $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOBIN)/golangci-lint

.PHONY: all test lint fmt vet tools

all: fmt vet lint test

test:
	@echo "Ejecutando pruebas..."
	@go test -v ./...

lint: tools
	@echo "Ejecutando linter..."
	@$(GOLANGCI_LINT) run

fmt:
	@echo "Formateando código..."
	@go fmt ./...

vet:
	@echo "Verificando código..."
	@go vet ./...

tools:
	@echo "Instalando herramientas necesarias..."
	@if [ ! -f $(GOLANGCI_LINT) ]; then \
		echo "golangci-lint no encontrado, instalando..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Herramientas instaladas correctamente."
