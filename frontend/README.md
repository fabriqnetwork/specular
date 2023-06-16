# xDAI to ETH 

Simple tool to bridge xDAI tokens to ETH on SPecular Network

## Getting Started

Clone the repo:


Move into the project directory:

```sh
cd specularBridge
```

Install project dependencies:

```sh
npm install
```

Create the required `.env` file from the example provided in the repo:

- [Chiado Testnet](./.env.chiado)
```sh
cp .env.chiado .env
```

Run the Dapp in development mode. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits. You will also see any lint errors in the console.

```sh
npm start
```

Run the Dapp in production mode. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

```sh
npm run build     
serve -s build -p 3000
```

Run the Dapp in Docker Container. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

```sh
docker-compose up -d --build
```
