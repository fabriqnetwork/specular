#!/bin/sh
cd "$(dirname "$0")" && cd ..

npx ts-node create_genesis.ts --in base_genesis.json --out ../clients/geth/specular/data/genesis.json
