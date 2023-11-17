# Ops

Specular chain operation tools.

## Bindings

```bash
make -C bindings
```

## Genesis generation

```bash
go run ./cmd/genesis/main.go \
    --genesis-config ./genesis-config.json \
    --out ./genesis.json \
    --l1-rpc-url http://localhost:8545 \
    --l1-block 0 \
    --export-hash ./genesis_hash.json
```
