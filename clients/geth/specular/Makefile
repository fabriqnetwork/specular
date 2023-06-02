.PHONY: clean install geth-docker

GOBIN = ./build/bin
# geth
GETH_SRC = ./cmd/geth/
GETH_TARGET = $(GOBIN)/geth
# clef
CLEF_SRC = ../go-ethereum/cmd/clef/
CLEF_TARGET = $(GOBIN)/clef
# bindings
BINDINGS_SRC = ./...
BINDINGS_TARGET = ./bindings
# contracts; TODO: use env variable
CONTRACTS_DIR = ../../../contracts
CONTRACTS_SRC = $(CONTRACTS_DIR)/src
CONTRACTS_TARGET = $(CONTRACTS_DIR)/artifacts/build-info

# PHONY targets

# Install all targets
install: $(GETH_TARGET) $(CLEF_TARGET)
# For back-compat (TODO: remove?)
geth-docker: $(GETH_TARGET)
# Remove:
# - bindings (do not remove bindings/gen.go)
# - contracts (this has to happen after bindings)
# - geth and clef
clean:
	rm $(BINDINGS_TARGET)/I*.go
	cd $(CONTRACTS_DIR) && npx hardhat clean
	rm -r $(GETH_TARGET)
	rm -r $(CLEF_TARGET)

# Real targets

# prereqs: all new/deleted files in contracts/ AND existing solidity files
$(CONTRACTS_TARGET): $(CONTRACTS_SRC) $(shell find $(CONTRACTS_DIR) -type f -name "*.sol")
	./sbin/compile_contracts.sh

# `touch` ensures the target is newer than preqreqs.
# This is required since `go generate` may not add/delete files.
$(BINDINGS_TARGET): $(CONTRACTS_TARGET)
	go generate $(BINDINGS_SRC)
	touch $(BINDINGS_TARGET)

$(GETH_TARGET): $(BINDINGS_TARGET)
	go build -o $(GETH_TARGET) $(GETH_SRC)
	@echo "Done building geth."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

$(CLEF_TARGET): $(CLEF_SRC)
	go build -o $(CLEF_TARGET) $(CLEF_SRC)
	@echo "Done building clef."
	@echo "Run \"$(GOBIN)/clef\" to launch clef."
