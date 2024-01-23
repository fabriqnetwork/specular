FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest as build

FROM golang:bullseye

ENV NODE_MAJOR=16
ENV FOUNDRY_VERSION="nightly"
ENV FOUNDRY_TAR="foundry_nightly_linux_amd64.tar.gz"

WORKDIR /tmp

RUN apt install -y python3 ca-certificates curl gnupg
RUN mkdir -p /etc/apt/keyrings && \
        curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
        echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
        apt-get update && \
        apt-get install nodejs -y

RUN wget https://github.com/foundry-rs/foundry/releases/download/$FOUNDRY_VERSION/$FOUNDRY_TAR && \
    tar xzvf $FOUNDRY_TAR && \
    mv cast /usr/local/bin

# copy everything we need

RUN mkdir -p /specular/workspace
RUN mkdir -p /specular/sbin
RUN mkdir -p /specular/contracts

# install hardhar
COPY --from=build /specular/package.json /specular
COPY --from=build /specular/contracts /specular/contracts
COPY --from=build /specular/pnpm-lock.yaml /specular/pnpm-lock.yaml
COPY --from=build /specular/pnpm-workspace.yaml /specular/pnpm-workspace.yaml

WORKDIR /specular

RUN npm install -g pnpm
RUN pnpm install


COPY --from=build /specular/config/local_docker /specular/workspace
COPY --from=build /specular/sbin /specular/sbin

# DEBUG/LOCAL BUILD
COPY ../config/local_docker /specular/workspace
COPY ../sbin /specular/sbin
# COPY ../services /specular/services

WORKDIR /specular/workspace

COPY --from=build /specular/ops/ /specular/ops/
COPY --from=build /specular/services/sidecar/build/bin/sidecar  /usr/local/bin/sidecar
COPY --from=build /specular/ops/build/bin/genesis  /usr/local/bin/genesis
COPY --from=build /specular/services/cl_clients/magi/target/debug/magi /usr/local/bin/magi
COPY --from=build /specular/services/el_clients/go-ethereum/build/bin/geth /usr/local/bin/geth


EXPOSE 4011 4012 4013 8545

