FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest as build

FROM alpine:3.19


COPY --from=build /specular/services/sidecar/build/bin/sidecar  /usr/local/bin/sidecar
COPY --from=build /specular/services/cl_clients/magi/target/debug/magi /usr/local/bin/magi
COPY --from=build /specular/services/el_clients/go-ethereum/build/bin/geth /usr/local/bin/geth

EXPOSE 4011 4012 4013

ENTRYPOINT ["/sbin/start-spgeth.sh"]
