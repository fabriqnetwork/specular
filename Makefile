.PHONY: install specular clean geth-docker contracts

SIDECAR_DIR = services/sidecar
SIDECAR_BIN = $(SIDECAR_DIR)/build/bin

CONTRACTS_DIR = contracts/
CONTRACTS_SRC = $(CONTRACTS_DIR)/src
CONTRACTS_TARGET = $(CONTRACTS_DIR)/artifacts/build-info

BINDINGS_TARGET = ./bindings

GETH_SRC = $(SIDECAR_DIR)/cmd/geth/
GETH_TARGET = $(SIDECAR_BIN)/geth

CLEF_SRC = $(GETH_SRC)/go-ethereum/cmd/clef/
CLEF_TARGET = $(GOBIN)/clef

install: sidecar $(GETH_TARGET) $(CLEF_TARGET)

geth: $(GETH_TARGET)

sidecar: $(BINDINGS_TARGET) $(shell find $(SIDECAR_DIR) -type f -name "*.go")
	cd $(SIDECAR_DIR)
	go build ./...

contracts: $(CONTRACTS_TARGET) # for back-compat

# Removes:
# - bindings (do not remove bindings/gen.go)
# - contracts (this has to happen after bindings)
# - geth and clef

# Docker process skips geth prereqs for docker building.
geth-docker: bindings-docker
	go build -o $(GETH_TARGET) $(GETH_SRC)
	@echo "Done building geth."
	@echo "Run \"$(GETH_TARGET)\" to launch geth."

bindings-docker:
	go generate ./...
	touch $(BINDINGS_TARGET)

# prereqs: all new/deleted files in contracts/ AND existing solidity files
$(CONTRACTS_TARGET): $(CONTRACTS_SRC) $(shell find $(CONTRACTS_DIR) -type f -name "*.sol")
	./sbin/compile_contracts.sh

# `touch` ensures the target is newer than preqreqs.
# This is required since `go generate` may not add/delete files.
$(BINDINGS_TARGET): $(CONTRACTS_TARGET)
	go generate ./...
	touch $(BINDINGS_TARGET)

$(GETH_TARGET): $(BINDINGS_TARGET)
	go build -o $(GETH_TARGET) $(GETH_SRC)
	@echo "Done building geth."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

$(CLEF_TARGET): $(CLEF_SRC)
	go build -o $(CLEF_TARGET) $(CLEF_SRC)
	@echo "Done building clef."
	@echo "Run \"$(GOBIN)/clef\" to launch clef."

clean:
	rm -f $(BINDINGS_TARGET)/I*.go
	cd $(CONTRACTS_DIR) && npx hardhat clean
	rm -rf $(GETH_TARGET)
	rm -rf $(CLEF_TARGET)
