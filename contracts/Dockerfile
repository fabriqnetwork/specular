FROM node:16-alpine3.17
COPY . /specular

RUN apk add --no-cache git python3 make g++

WORKDIR /specular/contracts

RUN yarn install --frozen-lockfile && npx hardhat compile

EXPOSE 8545
CMD ["npx", "hardhat", "node", "--hostname", "0.0.0.0"]
