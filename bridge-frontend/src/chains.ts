import {
  CHIADO_EXPLORER_URL,
  CHIADO_RPC_URL,
  SPECULAR_EXPLORER_URL,
  SPECULAR_RPC_URL,
  BRIDGE_GAS_LIMIT,
} from "./constants";

interface NativeCurrency {
  name: string;
  symbol: string;
  decimals: number;
}

interface Network {
  forkVersion: string;
  name: string;
  symbol: string;
  chainName: string;
  rpcUrl: string;
  blockExplorerUrl: string;
  nativeCurrency: NativeCurrency;
  chainId: string; // Add the chainId property here
}

const NETWORKS: Record<string, Network> = {
  '31337': {
    forkVersion: "0000006f",
    name: 'Chiado',
    symbol: 'XDAI',
    chainName: 'Chiado Testnet',
    rpcUrl: CHIADO_RPC_URL,
    blockExplorerUrl: CHIADO_EXPLORER_URL,
    nativeCurrency: {
      name: "xDAI",
      symbol: "xDAI",
      decimals: 18,
    },
    chainId: '31337', // Add the chainId value here
  },
  '13527': {
    forkVersion: "00000064",
    name: 'Specular Devnet',
    symbol: 'ETH',
    chainName: 'specular devnet',
    rpcUrl: SPECULAR_RPC_URL,
    blockExplorerUrl: SPECULAR_EXPLORER_URL,
    nativeCurrency: {
      name: "ETH",
      symbol: "ETH",
      decimals: 18,
    },
    chainId: '93481', // Add the chainId value here
  }
};

export {
  BRIDGE_GAS_LIMIT,
  NETWORKS
};
