import { ethers } from 'ethers'
import {
    L1ChainID,
    L2ChainID,
    L1Contracts,
    L2Contracts,
    ContractsLike
} from '../interfaces'

// L1_RPC_URL = http://l1-geth:8545
// L2_RPC_URL = http://sp-geth:4011

const l1portalAddresses = {
    sepolia: '0x457818a300BBa8E998889FAb2e8453BDebA238E0',
    hardhat_local: '0x6FfcAf7eDC12164d5f4651674f2957842f50a548'
}
const l1StandardBridgeAddresses = {
    sepolia: '0x06B8bF239246Cc0B08B77E743d816A690C7b1325',
    hardhat_local: '0xF60BE49fe797Dec1e3d8A084C86cb500211e74a7'
}
const l1RollupAddresses = {
    sepolia: '0x849eB65348360a6a5f33Bce19249763fFc426D5C',
    hardhat_local: '0x9D0B223A3e37F053C610a024b452676A82F558f2'
}


export const DEPOSIT_CONFIRMATION_BLOCKS: {
    [ChainID in L2ChainID]: number
} = {
    [L2ChainID.SPECULAR]: 2 as const,
    [L2ChainID.SPECULAR_HARDHAT_LOCAL]: 1 as const,
}

export const CHAIN_BLOCK_TIMES: {
    [ChainID in L1ChainID]: number
} = {
    [L1ChainID.MAINNET]: 13 as const,
    [L1ChainID.SEPOLIA]: 15 as const,
    [L1ChainID.HARDHAT_LOCAL]: 1 as const,
}


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
 * Mapping of L1 chain IDs to the appropriate contract addresses for the 
 * given network. Simplifies the process of getting the correct contract addresses for a given
 * contract name.
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
