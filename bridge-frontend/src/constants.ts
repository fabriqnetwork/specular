import { ethers } from "ethers";

export const CHIADO_NETWORK_ID = 31337;
export const CHIADO_RPC_URL = "http://localhost:8545";
export const CHIADO_EXPLORER_URL = "https://blockscout.chiadochain.net";
export const SPECULAR_NETWORK_ID = 13527;
export const SPECULAR_RPC_URL = "http://localhost:4011";
export const SPECULAR_EXPLORER_URL = "https://explorer.specular.network";

export const L1PORTAL_ADDRESS = "0x13D69Cf7d6CE4218F646B759Dcf334D82c023d8e";
export const L2PORTAL_ADDRESS = "0xBC9129Dc0487fc2E169941C75aABC539f208fb01";
export const L1ORACLE_ADDRESS = "0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1";
export const INBOX_ADDRESS = "0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1";
export const ROLLUP_ADDRESS = "0xF6168876932289D073567f347121A267095f3DD6";

export const BRIDGE_GAS_LIMIT = 300000
export const DEPOSIT_BALANCE_THRESHOLD = ethers.utils.parseEther("0.009");
export const RELAYER_ENDPOINT = "https://api.specular.network/devnet";

export const L1_BRIDGE_ADDR = "0xE7C2a73131dd48D8AC46dCD7Ab80C8cbeE5b410A";
export const L2_BRIDGE_ADDR = "0xF6168876932289D073567f347121A267095f3DD6";
