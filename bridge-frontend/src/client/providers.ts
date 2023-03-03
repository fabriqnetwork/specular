import { ethers } from "ethers";
import type { ExternalProvider } from "@ethersproject/providers";
import { jsonRpcProvider } from "vagmi/providers/jsonRpc";
import { CHIADO_RPC_URL, SPECULAR_RPC_URL } from "./constants";

const chiadoProvider = {
  http: CHIADO_RPC_URL,
};

const specularProvider = {
  http: SPECULAR_RPC_URL,
};

export const provider = jsonRpcProvider({
  rpc: (chain) => {
    if (chain.name === "Chiado") {
      return chiadoProvider;
    } else if (chain.name === "Specular Devnet") {
      return specularProvider;
    } else {
      return null;
    }
  },
});

export const transactor = new ethers.providers.Web3Provider(
  window.ethereum as ExternalProvider
);
