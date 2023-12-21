FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest as build

FROM alpine:3.19

COPY --from=build /specular/services/el_clients/go-ethereum/build/bin/geth /sbin/sp-geth
