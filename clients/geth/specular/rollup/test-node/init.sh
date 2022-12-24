#!/bin/bash
../../build/bin/geth --datadir ./data --networkid 13527 init ./genesis.json
../../build/bin/geth --datadir ./data_validator --networkid 13527 init ./genesis.json
../../build/bin/geth --datadir ./data_indexer --networkid 13527 init ./genesis.json
