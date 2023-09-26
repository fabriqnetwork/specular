#!/bin/sh
cd "$(dirname "$0")" && cd ..

pwd
npx ts-node src/create_genesis.ts --in data/base_genesis.json --out ../services/sidecar/data/genesis.json
