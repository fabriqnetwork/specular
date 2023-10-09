.PHONY: install sidecar clean geth-docker contracts

SIDECAR_DIR = services/sidecar
SIDECAR_BIN_SRC = ./cmd/sidecar
SIDECAR_BIN_TARGET = ./build/bin/sidecar
SIDECAR_BINDINGS_TARGET = $(SIDECAR_DIR)/bindings

CONTRACTS_DIR = contracts/
CONTRACTS_SRC = $(CONTRACTS_DIR)/src
CONTRACTS_TARGET = $(CONTRACTS_DIR)/artifacts/build-info

#TODO migrate to services/el_clients/go-ethereum
GETH_SRC = services/el_clients/go-ethereum
GETH_BIN_SRC = ./cmd/geth/
GETH_BIN_TARGET = ./build/bin/geth

#CLEF_SRC = $(SIDECAR_DIR)/cmd/clef/
#CLEF_TARGET = $(SIDECAR_BIN)/clef

# TODO add clef back in when moving to services/el_clients/go-ethereum
# install: sidecar $(GETH_TARGET) $(CLEF_TARGET)
install: $(GETH_BIN_TARGET)

geth: $(GETH_BIN_TARGET)

sidecar: $(SIDECAR_BINDINGS_TARGET) $(shell find $(SIDECAR_DIR) -type f -name "*.go")
	cd $(SIDECAR_DIR) && go build -o $(SIDECAR_BIN_TARGET) $(SIDECAR_BIN_SRC)

contracts: $(CONTRACTS_TARGET) # for back-compat

clean:
	rm -f $(SIDECAR_BINDINGS_TARGET)/I*.go
	cd $(CONTRACTS_DIR) && npx hardhat clean
	rm -rf $(SIDECAR_BIN_TARGET)
	rm -rf $(GETH_BIN_TARGET)
	#rm -rf $(CLEF_TARGET)

# Removes:
# - bindings (do not remove bindings/gen.go)
# - contracts (this has to happen after bindings)
# - geth and clef

# Docker process skips geth prereqs for docker building.
geth-docker: bindings-docker
	go build -o ./$(GETH_BIN_TARGET) ./$(GETH_BIN_SRC)
	@echo "Done building geth."
	@echo "Run \"$(GETH_BIN_TARGET)\" to launch geth."

bindings-docker:
	cd $(SIDECAR_DIR) && go generate ./...
	touch $(SIDECAR_BINDINGS_TARGET)

# prereqs: all new/deleted files in contracts/ AND existing solidity files
$(CONTRACTS_TARGET): $(CONTRACTS_SRC) $(shell find $(CONTRACTS_DIR) -type f -name "*.sol")
	sbin/compile_contracts.sh

# `touch` ensures the target is newer than preqreqs.
# This is required since `go generate` may not add/delete files.
$(SIDECAR_BINDINGS_TARGET): $(CONTRACTS_TARGET)
	cd $(SIDECAR_DIR) && go generate ./...
	touch $(SIDECAR_BINDINGS_TARGET)

$(GETH_BIN_TARGET): $(SIDECAR_BINDINGS_TARGET)
	cd $(GETH_SRC) && go build -o $(GETH_BIN_TARGET) $(GETH_BIN_SRC)
	@echo "Done building geth."
	#@echo "Run \"$(GOBIN)/geth\" to launch geth."

#$(CLEF_TARGET): $(CLEF_SRC)
	#go build -o ./$(CLEF_TARGET) ./$(CLEF_SRC)
	#@echo "Done building clef."
	##@echo "Run \"$(GOBIN)/clef\" to launch clef."
