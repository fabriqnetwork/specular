FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest as build

FROM golang:bullseye

ENV NODE_MAJOR=20

RUN apt install -y python3 ca-certificates curl gnupg
RUN mkdir -p /etc/apt/keyrings && \
        curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
        echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
        apt-get update && \
        apt-get install nodejs -y

RUN mkdir -p /specular/workspace
RUN mkdir -p /specular/sbin
RUN mkdir -p /specular/contracts

WORKDIR /specular/workspace

COPY --from=build /specular/config/local_docker/.contracts.env /specular/config/local_docker/.genesis.env /specular/config/local_docker/.paths.env /specular/config/local_docker/.sidecar.env /specular/config/local_docker/.sp_geth.env /specular/config/local_docker/.sp_magi.env /specular/config/local_docker/base_sp_rollup.json /specular/config/local_docker/genesis_config.json /specular/workspace

RUN cp /specular/workspace/base_sp_rollup.json /specular/workspace/sp_rollup.json

COPY --from=build /specular/sbin/ /specular/sbin/

COPY --from=build /specular/services/sidecar/build/bin/sidecar  /usr/local/bin/sidecar
COPY --from=build /specular/services/cl_clients/magi/target/debug/magi /usr/local/bin/magi
COPY --from=build /specular/services/el_clients/go-ethereum/build/bin/geth /usr/local/bin/geth


EXPOSE 4011 4012 4013 8545

