SHELL := /usr/bin/env bash

.PHONY: all
all: install-pre-commit-mac \
	install-pre-commit-linux \
	install-pre-commit-hooks \
	golangci \
	yamllint \
	govet \
	goerrcheck \
	gofull \
	test-report \
	test \
	lint


install-pre-commit-mac:
	@brew install pre-commit
install-pre-commit-linux:
	@sudo apt install pre-commit
install-pre-commit-hooks:
	@pre-commit install --install-hooks
	@pre-commit install --hook-type commit-msg --install-hooks

gofmt:
	@go fmt ./...


golangci:
	@golangci-lint run ./...

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

.PHONY: gost-report
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

.PHONY: lint
test: ## Run all tests
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

lint: ## Run all linters
	@echo "Running linters..."
	make yamllint
	@$(MAKE) golangci
	@$(MAKE) goerrcheck
	@$(MAKE) govet
