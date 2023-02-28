.PHONY: geth bindings install

GOBIN = ./build/bin

install: bindings geth

contracts:
	./sbin/compile_contracts.sh

bindings: contracts
	go generate ./...

geth: bindings
	go build -o $(GOBIN)/geth ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

geth-docker: 
	go build -o $(GOBIN)/geth ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."