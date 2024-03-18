import { ethers } from "ethers";
import * as dotenv from "dotenv";
dotenv.config({ path: ".env.test" });

import {
  L1ChainID,
  L2ChainID,
  L1Contracts,
  L2Contracts,
  ContractsLike,
} from "../interfaces";

// Docker
// L1_RPC_URL = http://l1-geth:8545
// L2_RPC_URL = http://sp-geth:4011

const l1portalAddresses = {
  sepolia: "0x457818a300BBa8E998889FAb2e8453BDebA238E0",
  hardhat_local: "0xDD82cA583E76aeDE3b050ab01C4aE3A834Eef683",
};
const l1StandardBridgeAddresses = {
  sepolia: "0x06B8bF239246Cc0B08B77E743d816A690C7b1325",
  hardhat_local: "0x94fCb90624FEf50984077A93808859B16bcEF44A",
};
const l1RollupAddresses = {
  sepolia: "0x849eB65348360a6a5f33Bce19249763fFc426D5C",
  hardhat_local: "0xC1DE463BA76396eCB27A8392bbdFEf6c641b2755",
};

/**
 * Full list of default L2 contract addresses.
 */
export const DEFAULT_L2_CONTRACT_ADDRESSES: L2Contracts = {
  UUPSPlaceholder: "0x2A00000000000000000000000000000000000000",
  L1Oracle: "0x2A00000000000000000000000000000000000010",
  L2Portal: "0x2A00000000000000000000000000000000000011",
  L2StandardBridge: "0x2A00000000000000000000000000000000000012",
  L1FeeVault: "0x2A00000000000000000000000000000000000020",
  L2BaseFeeVault: "0x2A00000000000000000000000000000000000021",
};

/**
 * Loads the L1 contracts for a given network by the network Id.
 *
 * @param network The ID of the network to load the contracts for.
 * @returns The L1 contracts for the given network.
 */
export const getL1ContractsByNetworkId = (network: number): L1Contracts => {
  if (network == L1ChainID.SEPOLIA) {
    return {
      L1Portal: l1portalAddresses.sepolia,
      L1StandardBridge: l1StandardBridgeAddresses.sepolia,
      L1Rollup: l1RollupAddresses.sepolia,
    };
  } else {
    return {
      L1Portal: l1portalAddresses.hardhat_local,
      L1StandardBridge: l1StandardBridgeAddresses.hardhat_local,
      L1Rollup: l1RollupAddresses.hardhat_local,
    };
  }
};

/**
 * Maps L1 chain IDs to the appropriate contract addresses of the given network.
 */
export const CONTRACT_ADDRESSES: {
  [ChainID in L2ChainID]: ContractsLike;
} = {
  [L2ChainID.SPECULAR]: {
    l1: getL1ContractsByNetworkId(L1ChainID.SEPOLIA),
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
  [L2ChainID.SPECULAR_HARDHAT_LOCAL]: {
    l1: getL1ContractsByNetworkId(L1ChainID.HARDHAT_LOCAL),
    l2: DEFAULT_L2_CONTRACT_ADDRESSES,
  },
};
