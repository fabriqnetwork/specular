SHELL := /usr/bin/env bash

pkg := bindings
monorepo-base := $(shell dirname $(realpath ..))
contracts-dir := $(monorepo-base)/contracts
geth-dir := $(monorepo-base)/services/el_clients/go-ethereum

all: bindings

abigen:
	cd $(geth-dir) && \
		go run build/ci.go install ./cmd/abigen

compile:
	cd $(contracts-dir) && \
		pnpm build

bindings: compile abigen
	go run ./gen/main.go \
		-out ./bindings \
		-contracts ./artifacts.json \
		-package $(pkg) \
		-contract-dir $(contracts-dir) \
		-abigen-bin $(geth-dir)/build/bin/abigen