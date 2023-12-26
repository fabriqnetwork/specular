FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:specular-latest

WORKDIR /tmp
ENV FOUNDRY_VERSION="nightly-67ab8704476d55e47545cf6217e236553c427a80"
ENV FOUNDRY_TAR="foundry_nightly_linux_amd64.tar.gz"
RUN wget https://github.com/foundry-rs/foundry/releases/download/$FOUNDRY_VERSION/$FOUNDRY_TAR && \
    tar xzvf $FOUNDRY_TAR && \
    mv cast /usr/local/bin
WORKDIR /specular/workspace
