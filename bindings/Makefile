#!/bin/bash

GO_CMD=go
GO_BUILD=$(GO_CMD) build
CMD_PATH=./gen/main.go
DIST=dist
BINARY_NAME=bindings

PKG="bindings"
MONOREPO_BASE=`git rev-parse --show-toplevel | xargs`
CONTRACTS_DIR="${MONOREPO_BASE}/contracts"

.PHONY: all # Run all necessary steps and generate bindings
all: bindings

.PHONY: abigen # Install latest abigen command binary
abigen:
	@go install github.com/ethereum/go-ethereum/cmd/abigen@latest

.PHONY: compile # Compile contracts from the monorepo
compile:
	@echo "monorepo-base: ${MONOREPO_BASE}"
	@echo "contracts-dir: ${CONTRACTS_DIR}"
	@cd "${CONTRACTS_DIR}" && \
		pnpm build

.PHONY: build # Build binary
build:
	mkdir -p $(DIST)
	CGO_ENABLED=1 $(GO_BUILD) -o ./dist/$(BINARY_NAME) -v $(CMD_PATH)

.PHONY: test # Run tests
test:
	go test $(go list ./... | grep -v /bindings/bindings)

.PHONY: bindings # Run bindings generation alone
bindings: compile abigen
	@go run ./gen/main.go \
		-out ./bindings \
		-contracts ./artifacts.json \
		-package "${PKG}" \
		-contract-dir "${CONTRACTS_DIR}" \
		-abigen-bin abigen

.PHONY: lint-test # Run lint-tests
lint-test:
	go mod tidy
	@go install golang.org/x/tools/cmd/goimports@latest
	goimports -local github.com/SpecularL2/bindings -w .
	go fmt ./...
	@golangci-lint -v run ./...
	@test -z "$$(golangci-lint run ./...)"

.PHONY: lint # Run lint-tests (alias)
lint: lint-test

.PHONY: help # Help - list of targets with descriptions
help:
	@echo ''
	@echo 'Usage:'
	@echo ' ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/ \1\t\2/' | expand -t20
