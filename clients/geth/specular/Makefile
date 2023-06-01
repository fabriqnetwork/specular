.PHONY: geth bindings install

GOBIN = ./build/bin
GETH_ENTRY = ./entry/geth/cmd/geth

install: bindings geth clef

contracts:
	./sbin/compile_contracts.sh

bindings: contracts
	go generate ./...

geth: bindings
	go build -o $(GOBIN)/geth $(GETH_ENTRY)
	@echo "Done building geth."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

clef:
	go build -o $(GOBIN)/clef ./../go-ethereum/cmd/clef
	@echo "Done building clef."
	@echo "Run \"$(GOBIN)/clef\" to launch clef."

geth-docker:
	go build -o $(GOBIN)/geth $(GETH_ENTRY)
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."
