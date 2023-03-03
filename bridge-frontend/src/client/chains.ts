import type { Chain } from "vagmi";

import {
  CHIADO_EXPLORER_URL,
  CHIADO_RPC_URL,
  SPECULAR_EXPLORER_URL,
  SPECULAR_RPC_URL,
} from "./constants";

export const chiadoChain: Chain = {
  id: 10200,
  name: "Chiado",
  network: "chiado",
  nativeCurrency: {
    name: "xDAI",
    symbol: "xDAI",
    decimals: 18,
  },
  rpcUrls: {
    default: CHIADO_RPC_URL,
  },
  blockExplorers: {
    default: {
      name: "Chiado Explorer",
      url: CHIADO_EXPLORER_URL,
    },
  },
  testnet: true,
};

export const specularChain: Chain = {
  id: 93481,
  name: "Specular Devnet",
  network: "specular devnet",
  nativeCurrency: {
    name: "ETH",
    symbol: "ETH",
    decimals: 18,
  },
  rpcUrls: {
    default: SPECULAR_RPC_URL,
  },
  blockExplorers: {
    default: {
      name: "Specular Explorer",
      url: SPECULAR_EXPLORER_URL,
    },
  },
  testnet: true,
};
