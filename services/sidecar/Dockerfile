# install dependencies and build sidecar binary
FROM node:16-alpine3.17 as builder

RUN apk add --no-cache git go python3 make g++ musl-dev linux-headers
RUN corepack enable

WORKDIR /specular
COPY . /specular

# frozen lockfile is automatically enabled in CI environments
RUN pnpm install
RUN make sidecar

# create container with just the binary
FROM alpine:3.17

WORKDIR /specular

RUN apk add --no-cache ca-certificates bash
COPY --from=builder /specular/services/sidecar/build/bin/sidecar /usr/local/bin/

EXPOSE 8545 8546
