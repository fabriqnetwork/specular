#FROM ethereum/client-go:alltools-latest-amd64

# RUN apk add --no-cache bash nodejs-current npm python3 make g++ go musl-dev linux-headers git curl
# RUN corepack enable
# RUN apk add --no-cache alpine-sdk openssl-dev pkgconfig
# RUN curl https://sh.rustup.rs -o rustup-init.sh && sh rustup-init.sh -y
# ENV PATH="${PATH}:/root/.cargo/bin"

FROM --platform=$BUILDPLATFORM debian:bullseye-slim as base

# Versions
ENV GOVERSION=1.21.5
ENV NODE_MAJOR=20

RUN apt update && apt install -y ca-certificates curl gnupg
RUN mkdir -p /etc/apt/keyrings
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg

# Node Version
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list

RUN apt-get update && apt-get install nodejs -y

RUN apt install -y libudev-dev build-essential ca-certificates clang libpq-dev libssl-dev pkg-config lsof lld bash python3 make g++ musl-dev git

RUN apt install -y linux-headers-generic wget
RUN npm install -g pnpm
RUN wget https://go.dev/dl/go$GOVERSION.linux-amd64.tar.gz && tar -C /usr/local -xzf go$GOVERSION.linux-amd64.tar.gz && rm go$GOVERSION.linux-amd64.tar.gz


RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH "$PATH:/root/.cargo/bin"
ENV PATH="$PATH:/usr/local/go/bin"


