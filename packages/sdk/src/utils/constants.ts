import { ethers } from 'ethers'
import {
    L1ChainID,
    L2ChainID,
    L1Contracts,
    L2Contracts,
    ContractsLike
} from '../interfaces'

// Docker
// L1_RPC_URL = http://l1-geth:8545
// L2_RPC_URL = http://sp-geth:4011

const l1portalAddresses = {
    sepolia: '0x457818a300BBa8E998889FAb2e8453BDebA238E0',
    hardhat_local: '0x42bB391bd8A51f024E1D59726888186b15E08359'
}
const l1StandardBridgeAddresses = {
    sepolia: '0x06B8bF239246Cc0B08B77E743d816A690C7b1325',
    hardhat_local: '0x6eE4BBa799CB77004CD3f21AE074f061194b2C20'
}
const l1RollupAddresses = {
    sepolia: '0x849eB65348360a6a5f33Bce19249763fFc426D5C',
    hardhat_local: '0xD2BaAdF3f3E5e8868F8aEE6E564D1ceB56461040'
}

export const l1ChainIds = [1337, 11155111]

export const l2ChainIds = [13527]


/**
 * Full list of default L2 contract addresses.
 */
export const DEFAULT_L2_CONTRACT_ADDRESSES: L2Contracts = {
    UUPSPlaceholder: "0x2A00000000000000000000000000000000000000",
    L1Oracle: "0x2A00000000000000000000000000000000000010",
    L2Portal: "0x2A00000000000000000000000000000000000011",
    L2StandardBridge: "0x2A00000000000000000000000000000000000012",
    L1FeeVault: "0x2A00000000000000000000000000000000000020",
    L2BaseFeeVault: "0x2A00000000000000000000000000000000000021"
}

/**
 * Loads the L1 contracts for a given network by the network Id.
 *
 * @param network The ID of the network to load the contracts for.
 * @returns The L1 contracts for the given network.
 */
export const getL1ContractsByNetworkId = (network: number): L1Contracts => {
    if (network == 11155111) {
        return {
            L1Portal: l1portalAddresses.sepolia,
            L1StandardBridge: l1StandardBridgeAddresses.sepolia,
            L1Rollup: l1RollupAddresses.sepolia,
        }
    } else {
        return {
            L1Portal: l1portalAddresses.hardhat_local,
            L1StandardBridge: l1StandardBridgeAddresses.hardhat_local,
            L1Rollup: l1RollupAddresses.hardhat_local,
        }
    }
}

/**
 * Maps L1 chain IDs to the appropriate contract addresses of the given network.
 */
export const CONTRACT_ADDRESSES: {
    [ChainID in L2ChainID]: ContractsLike
} = {
    [L2ChainID.SPECULAR]: {
        l1: getL1ContractsByNetworkId(11155111),
        l2: DEFAULT_L2_CONTRACT_ADDRESSES,
    },
    [L2ChainID.SPECULAR_HARDHAT_LOCAL]: {
        l1: getL1ContractsByNetworkId(31337),
        l2: DEFAULT_L2_CONTRACT_ADDRESSES,
    },
}
