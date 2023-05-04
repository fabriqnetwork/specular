FROM node:16-alpine3.17
COPY . /specular

RUN apk add --no-cache git

WORKDIR /specular/contracts
RUN apk add --no-cache python3 make g++
RUN yarn install --frozen-lockfile && npx hardhat compile

EXPOSE 8545
CMD ["npx", "hardhat", "node", "--hostname", "0.0.0.0"]
