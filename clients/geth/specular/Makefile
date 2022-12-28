.PHONY: geth bindings install

GOBIN = ./build/bin

install: bindings geth

bindings:
	go generate ./...

geth:
	go build -o $(GOBIN)/geth ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."
