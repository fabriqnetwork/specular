.PHONY: install sidecar clean geth-docker contracts ops

SIDECAR_DIR = services/sidecar
SIDECAR_BIN_SRC = ./cmd/sidecar # relative to SIDECAR_DIR
SIDECAR_BIN_TARGET = ./build/bin/sidecar # relative to SIDECAR_DIR

CONTRACTS_DIR = contracts/
CONTRACTS_SRC = $(CONTRACTS_DIR)/src
CONTRACTS_TARGET = $(CONTRACTS_DIR)/artifacts/build-info

GETH_SRC = services/el_clients/go-ethereum
GETH_BIN_TARGET = ./build/bin/geth # relative to GETH_SRC

MAGI_DIR = services/cl_clients/magi
MAGI_BIN_TARGET = services/cl_clients/magi/target/debug/magi

OPS_DIR = ops
OPS_BIN_TARGET = ./build/bin/genesis # relative to OPS_DIR
OPS_BIN_SRC = ./cmd/genesis/ # relative to OPS_DIR

BINDINGS_TARGET = ./bindings-go

# TODO add clef back in when moving to services/el_clients/go-ethereum
#CLEF_SRC = $(SIDECAR_DIR)/cmd/clef/
#CLEF_TARGET = $(SIDECAR_BIN)/clef

install: geth magi sidecar ops
geth: bindings $(GETH_BIN_TARGET)
magi: $(MAGI_BIN_TARGET)
sidecar: bindings $(shell find $(SIDECAR_DIR) -type f -name "*.go")
	cd $(SIDECAR_DIR) && go build -o $(SIDECAR_BIN_TARGET) $(SIDECAR_BIN_SRC)

bindings: $(CONTRACTS_TARGET)
	GOFLAGS="-buildvcs=false" make -C $(BINDINGS_TARGET)

ops: bindings
	cd $(OPS_DIR) && go build -o $(OPS_BIN_TARGET) $(OPS_BIN_SRC)

contracts: $(CONTRACTS_TARGET) # for back-compat

# Removes:
# - bindings (do not remove bindings/gen.go)
# - contracts (this has to happen after bindings)
# - geth and clef
clean:
	cd $(CONTRACTS_DIR) && npx hardhat clean
	rm -rf $(SIDECAR_BIN_TARGET)
	rm -rf $(GETH_BIN_TARGET)
	#rm -rf $(CLEF_TARGET)

# prereqs: all new/deleted files in contracts/ AND existing solidity files
$(CONTRACTS_TARGET): $(CONTRACTS_SRC) $(shell find $(CONTRACTS_DIR) -type f -name "*.sol")
	cd contracts && pnpm build

$(GETH_BIN_TARGET):
	cd $(GETH_SRC) && GOFLAGS="-buildvcs=false" $(MAKE) geth

$(MAGI_BIN_TARGET): $(shell find $(MAGI_DIR) -type f -name "*.rs")
	cd $(MAGI_DIR) && cargo build

#$(CLEF_TARGET): $(CLEF_SRC)
	#go build -o ./$(CLEF_TARGET) ./$(CLEF_SRC)
	#@echo "Done building clef."
	##@echo "Run \"$(GOBIN)/clef\" to launch clef."
