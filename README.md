# Specular Monorepo

## Directory Structure

<pre>
├── <a href="./services/">services</a>: Specular L2 clients
│   ├── <a href="./services/cl_clients">cl_clients</a>: Consensus Layer Clients
│   └── <a href="./services/el_clients/">el_clients</a>: Execution Layer Clients
│       └── <a href="./services/el_clients/geth/">bindings</a>: Minimally modified geth fork
├── <a href="./contracts">contracts</a>: Specular L1 and L2  contracts
└── <a href="./lib/">lib</a>: Libraries used in L2 EL Clients
│   └── <a href="./lib/el_golang_lib/">el_golang_lib</a>: Library for golang clients
</pre>

## License

Unless specified in subdirectories, this repository is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). See `LICENSE` for details.
