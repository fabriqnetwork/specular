#!/bin/sh
cd "$(dirname "$0")" && cd ..

npx ts-node src/create_genesis.ts --in data/base_genesis.json --out ../e2e/data/genesis.json

