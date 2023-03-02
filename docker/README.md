# Specular Docker Compose

## Prerequisites

### Create the project folder and add keys

_Note: `<specular>` refers to the path of the Specular monorepo._

```bash
mkdir <project> && cd <project>
mkdir specular-datadir && cd specular-datadir
cp <specular>/clients/geth/specular/data/keys/sequencer.prv ./key.prv # Change this key according to the configuration
cp <specular>/clients/geth/specular/data/password.txt .
```

### Generate `genesis.json`

```bash
cd <specular>/contracts
npx ts-node scripts/create_genesis.ts --in ../clients/geth/specular/data/base_genesis.json --out <project>/specular-datadir/genesis.json
```

## Run sequencer

```bash
cd <project>
docker compose -f <specular>/docker/docker-compose-sequencer.yml -p sequencer up -d
```

Sequencer will listen HTTP on port `4011` and WS on port `4012`. Note: if the sequencer container does not show up under `docker ps` then run `docker logs sequencer` to debug.

Remember to delete `geth` folder under `specular-datadir` on every restart.

## Run sequencer with block explorer

```bash
cd <project>
docker compose -f <specular>/docker/docker-compose-sequencer-explorer.yml -p sequencer-explorer up -d
```

Sequencer will listen HTTP on port `4011` and WS on port `4012`.
Blockscout is available on port `4000`.

## Run block explorer with indexer

```bash
cd <project>
docker compose -f <specular>/docker/docker-compose-explorer.yml -p explorer up -d
```

Blockscout is available on port `4000`.

## Run integration test environment

```bash
cd <project>
docker compose -f <specular>/docker/docker-compose-integration-tests.yml -p integration up -d
```

L1 hardhat node will listen HTTP on port `8545`.
Sequencer will listen HTTP on port `4011` and WS on port `4012`.
