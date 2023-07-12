import { ethers } from "ethers";

export const CHIADO_NETWORK_ID = 10200;
export const CHIADO_RPC_URL = "https://rpc.chiadochain.net";
export const CHIADO_EXPLORER_URL = "https://blockscout.chiadochain.net";
export const SPECULAR_NETWORK_ID = 93481;
export const SPECULAR_RPC_URL = "https://devnet.specular.network";
export const SPECULAR_EXPLORER_URL = "https://explorer.specular.network";
export const L1PORTAL_ADDRESS = "0x26b5fCaB7348a1B68827751ed03bcbe968484b58";
export const L2PORTAL_ADDRESS = "0x0A267A0C590c226FfeF9C7CF8DDE4a935Bf0864a";
export const L1ORACLE_ADDRESS = "0xb5776987319Ab81978e2F85568341b66E408466c";
export const INBOX_ADDRESS = "0x0F8Eb33De383D70EF7538113AB50c74a05AF4096";
export const ROLLUP_ADDRESS = "0x5Fdf6B833270A9562c735d2D0Cf784aC5b0fE8Cb";
export const BRIDGE_GAS_LIMIT = 300000
export const DEPOSIT_BALANCE_THRESHOLD = ethers.utils.parseEther("0.009");
export const RELAYER_ENDPOINT = "https://api.specular.network/devnet";
