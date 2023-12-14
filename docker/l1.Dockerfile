FROM specular:build-v0.0.1

# RUN apk add --no-cache bash nodejs-current npm python3 make g++ go musl-dev linux-headers git
# RUN corepack enable

WORKDIR /specular
COPY . /specular

# frozen lockfile is automatically enabled in CI environments
RUN pnpm install
RUN make

# TODO: what ports should be exposed?
EXPOSE 8545 8546
