# Specular Monorepo

## Directory Structure

<pre>
├── <a href="./clients/geth">clients/geth</a>: Specular L2 clients
│   ├── <a href="./clients/geth/go-ethereum">go-ethereum</a>: Minimally modified go-ethereum to support Specular prover
│   └── <a href="./clients/geth/specular">specular</a>: Specular client software
│       ├── <a href="./clients/geth/specular/bindings">bindings</a>: Golang bindings of Specular L1 contracts
│       ├── <a href="./clients/geth/specular/prover">proof</a>: Specular prover
│       └── <a href="./clients/geth/specular/rollup">rollup</a>: Specular rollup services
└── <a href="./contracts">contracts</a>: Specular L1 contracts
</pre>

## License

Unless specified in subdirectories, this repository is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). See `LICENSE` for details.
