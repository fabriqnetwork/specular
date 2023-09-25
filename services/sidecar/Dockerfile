# Compile contracts for binding generation
FROM node:16-alpine3.17 as contracts-builder

RUN apk add --no-cache bash git python3 make g++ musl-dev linux-headers && corepack enable

WORKDIR /specular
COPY pnpm-lock.yaml ./
RUN pnpm fetch

COPY . /specular
RUN pnpm install --offline --frozen-lockfile

WORKDIR /specular/services/sidecar
RUN make contracts

# Build geth binary
FROM golang:1.19-alpine3.17 as geth-builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /specular/clients/geth/go-ethereum
COPY --from=contracts-builder /specular/clients/geth/go-ethereum/go.mod ./
COPY --from=contracts-builder /specular/clients/geth/go-ethereum/go.sum ./
RUN go mod download

WORKDIR /specular/services/sidecar
COPY --from=contracts-builder /specular/services/sidecar/go.mod ./
COPY --from=contracts-builder /specular/services/sidecar/go.sum ./
RUN go mod download

COPY --from=contracts-builder /specular /specular

WORKDIR /specular/services/sidecar
RUN make geth-docker

# Specular geth image
FROM alpine:3.17

WORKDIR /specular

RUN apk add --no-cache ca-certificates
COPY --from=geth-builder /specular/services/sidecar/build/bin/geth /usr/local/bin/
COPY ./docker/scripts/geth.sh /usr/local/bin/

EXPOSE 8545 8546
ENTRYPOINT ["geth"]
