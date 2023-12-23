FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest as build

FROM scratch

COPY --from=build config/local_docker/.contracts.env config/local_docker/.genesis.env config/local_docker/.paths.env config/local_docker/.sidecar.env config/local_docker/.sp_geth.env config/local_docker/.sp_magi.env config/local_docker/base_sp_rollup.json config/local_docker/genesis_config.json

COPY --from=build ./sbin /sbin

COPY --from=build /specular/services/sidecar/build/bin/sidecar  /usr/local/bin/sidecar
COPY --from=build /specular/services/cl_clients/magi/target/debug/magi /usr/local/bin/magi
COPY --from=build /specular/services/el_clients/go-ethereum/build/bin/geth /usr/local/bin/geth

EXPOSE 4011 4012 4013

