# Compile contracts for binding generation
FROM node:16-alpine3.17 as contracts-builder

RUN apk add --no-cache make git python3 g++ corepack enable

COPY . /specular

WORKDIR /specular

RUN rm pnpm-lock.yaml
RUN rm pnpm-workspace.yaml

WORKDIR /specular/clients/geth/specular

RUN make contracts

# Build geth binary
FROM golang:1.19-alpine3.17 as geth-builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

COPY --from=contracts-builder /specular /specular

WORKDIR /specular/clients/geth/specular
RUN make geth-docker

# Specular geth image
FROM alpine:3.17

WORKDIR /specular

RUN apk add --no-cache ca-certificates
COPY --from=geth-builder /specular/clients/geth/specular/build/bin/geth /usr/local/bin/
COPY ./docker/scripts/geth.sh /usr/local/bin/

EXPOSE 8545 8546
ENTRYPOINT ["geth"]
