FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:build-latest

WORKDIR /tmp
RUN wget https://github.com/foundry-rs/foundry/releases/download/$FOUNDRY_VERSION/$FOUNDRY_TAR && \
    tar xzvf $FOUNDRY_TAR && \
    mv cast /usr/local/bin

WORKDIR /specular
ADD . /specular


# frozen lockfile is automatically enabled in CI environments
RUN pnpm install

ENV RUST_BACKTRACE=full
RUN make

# TODO: what ports should be exposed?
EXPOSE 8545 8546
