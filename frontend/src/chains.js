import {
  CHIADO_EXPLORER_URL,
  CHIADO_RPC_URL,
  SPECULAR_EXPLORER_URL,
  SPECULAR_RPC_URL,
  BRIDGE_GAS_LIMIT,
} from "./constants";

const NETWORKS = {
  10200: {
    forkVersion: "0000006f",
    name: 'Chiado',
    symbol: 'XDAI',
    chainName: 'Chiado Testnet',
    rpcUrl: CHIADO_RPC_URL,
    blockExplorerUrl: CHIADO_EXPLORER_URL
  },
  93481: {
    forkVersion: "00000064",
    name: 'Specular Devnet',
    symbol: 'ETH',
    chainName: 'specular devnet',
    rpcUrl: SPECULAR_RPC_URL,
    blockExplorerUrl: SPECULAR_EXPLORER_URL
  }
}

export {
  BRIDGE_GAS_LIMIT,
  NETWORKS
}