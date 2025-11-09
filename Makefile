SHELL := /usr/bin/env bash

# OS detection variable (Darwin for macOS, Linux for Linux)
UNAME_S := $(shell uname -s)

.PHONY: all install-pre-commit-mac install-pre-commit-linux install-pre-commit install-pre-commit-hooks gofmt golangci yamllint govet goerrcheck gofull gosec test-report test lint

all: install-pre-commit \
    install-pre-commit-hooks \
    golangci \
    yamllint \
    govet \
    goerrcheck \
    gofull \
    test-report \
    test \
    lint

# Runs only on macOS
install-pre-commit-mac:
	@echo "Checking pre-commit on macOS..."
	@if ! command -v pre-commit &> /dev/null; then \
		echo "pre-commit not found, installing via brew..."; \
		brew install pre-commit; \
	else \
		echo "pre-commit is already installed."; \
	fi

# Runs only on Linux
install-pre-commit-linux:
	@echo "Checking pre-commit on Linux..."
	@if ! command -v pre-commit &> /dev/null; then \
		echo "pre-commit not found, installing via apt..."; \
		sudo apt update && sudo apt install -y pre-commit; \
	else \
		echo "pre-commit is already installed."; \
	fi

# Selects the correct installer based on the operating system
install-pre-commit:
ifeq ($(UNAME_S), Darwin)
	@$(MAKE) install-pre-commit-mac
else ifeq ($(UNAME_S), Linux)
	@$(MAKE) install-pre-commit-linux
else
	@echo "Unsupported OS: $(UNAME_S). Please install pre-commit manually."
	@exit 1
endif

install-pre-commit-hooks:
	@pre-commit install --install-hooks
	@pre-commit install --hook-type commit-msg --install-hooks

gofmt:
	@go fmt ./...

golangci:
	@golangci-lint run --config=.golangci.yaml ./...

yamllint:
	@pre-commit run yamllint --all-files

govet:
	@go vet ./...

goerrcheck:
	@errcheck -ignore '[cC]lose' ./...

gofull:
	@golangci-lint run ./...
	@pre-commit run goerrcheck --all-files
	@pre-commit run go-vet --all-files

gosec:
	@gosec ./...

test-report:
	@echo "Running tests and generating coverage report..."
	@go test -coverprofile=coverage.out -covermode=atomic -race -v ./... || (echo "\nTests failed. Check test output above."; exit 1)
	@echo "\nGenerating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Generating function coverage report..."
	@go tool cover -func=coverage.out | tee coverage.txt
	@echo "\nCoverage Summary:"
	@go tool cover -func=coverage.out | grep total:
	@echo "\nHTML report: file://$(shell pwd)/coverage.html"

test: ## Run all tests
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

lint: ## Run all linters
	@echo "Running linters..."
	@$(MAKE) yamllint
	@$(MAKE) golangci
	@$(MAKE) goerrcheck
	@$(MAKE) govet
