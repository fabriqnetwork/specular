FROM specular:build-v0.0.1

WORKDIR /specular
ADD . /specular

# frozen lockfile is automatically enabled in CI environments
RUN pnpm install

ENV RUST_BACKTRACE=full
RUN make

# TODO: what ports should be exposed?
EXPOSE 8545 8546
