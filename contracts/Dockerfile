# Compile contracts for binding generation
FROM node:16-alpine3.17

RUN apk add --no-cache git python3 make g++ musl-dev linux-headers && corepack enable

COPY . /specular

WORKDIR /specular

RUN rm pnpm-lock.yaml
RUN rm pnpm-workspace.yaml

WORKDIR /specular/clients/geth/specular

RUN make contracts
RUN make geth

RUN ls sbin/
RUN chmod +x sbin/export_genesis.sh

WORKDIR /specular/contracts

RUN npx ts-node scripts/create_genesis.ts --in ../clients/geth/specular/data/base_genesis.json --out ../clients/geth/specular/data/genesis.json
RUN ls ../clients/geth/specular/sbin/

EXPOSE 8545
CMD ["npx", "hardhat", "node", "--hostname", "0.0.0.0"]
