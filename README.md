# Specular Monorepo

Welcome to the Specular monorepo, containing the code for the Specular L2 protocol and related code.
In this repo you will find:

<pre>
├── <a href="./clients/geth">clients/geth</a>: Specular L2 clients
│   ├── <a href="./clients/geth/go-ethereum">go-ethereum</a>: Minimally modified go-ethereum to support Specular prover
│   └── <a href="./clients/geth/specular">specular</a>: Specular client software
│       ├── <a href="./clients/geth/specular/bindings">bindings</a>: Golang bindings of Specular L1 contracts
│       ├── <a href="./clients/geth/specular/proof">proof</a>: Specular prover
│       └── <a href="./clients/geth/specular/rollup">rollup</a>: Specular rollup services
└── <a href="./contracts">contracts</a>: Specular L1 contracts
</pre>

## Getting Started for Developers

### Development Environment Requirements

To contribute to the Specular monorepo, it will be handy to have the Specular environment and
toolchain set up. The Specular toolchain depends on:

* go
* Node with pNPM

These can be installed on macOS with:

```
> brew install go pnpm
```

and on DEB-based systems:

```
> apt-get install go pnpm
```

### Submodule Dependencies

The Specular repository depends on a few other projects that are connected to the Specular
repository via git submodules including:

* blockscout
* go-ethereum
* forge-std

To collect and set up the required submodules, run:

```
> git submodule init
> git submodule update --recursive
```

### Building the Project Locally

TODO

## License

Unless specified in subdirectories, this repository is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). See `LICENSE` for details.
