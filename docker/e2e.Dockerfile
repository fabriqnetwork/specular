FROM 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:build-latest

WORKDIR /specular
ADD . /specular

# frozen lockfile is automatically enabled in CI environments
RUN pnpm install

ENV RUST_BACKTRACE=full
RUN make

# TODO: what ports should be exposed?
EXPOSE 8545 8546
