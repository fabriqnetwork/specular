.PHONY: install specular clean geth-docker contracts

GOSRC = ./...
GOBIN = ./build/bin
# TODO use env variables
SPECULAR_DIR = ./
CONTRACTS_DIR = ../../../contracts
# geth
GETH_SRC = ./cmd/geth/
GETH_TARGET = $(GOBIN)/geth
# clef
CLEF_SRC = ../go-ethereum/cmd/clef/
CLEF_TARGET = $(GOBIN)/clef
# bindings
BINDINGS_TARGET = ./bindings
# contracts
CONTRACTS_SRC = $(CONTRACTS_DIR)/src
CONTRACTS_TARGET = $(CONTRACTS_DIR)/artifacts/build-info

# PHONY targets

# Makes all targets
install: specular $(GETH_TARGET) $(CLEF_TARGET)
# Builds specular go source
# TODO: use real target to optimize build time
specular: $(BINDINGS_TARGET) $(shell find $(SPECULAR_DIR) -type f -name "*.go")
	go build $(GOSRC)
# Removes:
# - bindings (do not remove bindings/gen.go)
# - contracts (this has to happen after bindings)
# - geth and clef
clean:
	rm -f $(BINDINGS_TARGET)/I*.go
	cd $(CONTRACTS_DIR) && npx hardhat clean
	rm -rf $(GETH_TARGET)
	rm -rf $(CLEF_TARGET)
# Skips geth prereqs for docker building.
geth-docker:
	go build -o $(GETH_TARGET) $(GETH_SRC)
	@echo "Done building geth."
	@echo "Run \"$(GETH_TARGET)\" to launch geth."
contracts: $(CONTRACTS_TARGET) # for back-compat

# Real targets

# prereqs: all new/deleted files in contracts/ AND existing solidity files
$(CONTRACTS_TARGET): $(CONTRACTS_SRC) $(shell find $(CONTRACTS_DIR) -type f -name "*.sol")
	./sbin/compile_contracts.sh

# `touch` ensures the target is newer than preqreqs.
# This is required since `go generate` may not add/delete files.
$(BINDINGS_TARGET): $(CONTRACTS_TARGET)
	go generate $(GO_SRC)
	touch $(BINDINGS_TARGET)

$(GETH_TARGET): $(BINDINGS_TARGET)
	go build -o $(GETH_TARGET) $(GETH_SRC)
	@echo "Done building geth."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

$(CLEF_TARGET): $(CLEF_SRC)
	go build -o $(CLEF_TARGET) $(CLEF_SRC)
	@echo "Done building clef."
	@echo "Run \"$(GOBIN)/clef\" to launch clef."
