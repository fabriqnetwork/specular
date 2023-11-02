## Directory Structure

**TODO**

<pre>
├── <a href="./services/">services</a>: L2 services
│   ├── <a href="./services/cl_clients">cl_clients</a>: Consensus-layer clients
│   ├── <a href="./services/el_clients/">el_clients</a>: Execution-layer clients
│   │      └── <a href="./services/el_clients/go-ethereum/">go-ethereum</a>: Minimally modified geth fork
│   └── <a href="./services/sidecar/">sidecar</a>: Sidecar services
├── <a href="./contracts">contracts</a>: L1 and L2 contracts
└── <a href="./lib/">lib</a>: Libraries used in L2 EL Clients
    └── <a href="./lib/el_golang_lib/">el_golang_lib</a>: Library for golang EL clients
</pre>
