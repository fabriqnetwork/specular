FROM ethereum/client-go

RUN apk add --no-cache bash nodejs-current npm python3 make g++ go musl-dev linux-headers git curl
RUN corepack enable
RUN apk add --no-cache alpine-sdk openssl-dev pkgconfig
RUN curl https://sh.rustup.rs -o rustup-init.sh && sh rustup-init.sh -y
ENV PATH="${PATH}:/root/.cargo/bin"

