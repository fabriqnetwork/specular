import {
  Provider,
  TransactionReceipt,
  TransactionResponse,
} from '@ethersproject/abstract-provider'
import { Signer } from '@ethersproject/abstract-signer'
import { Contract, BigNumber } from 'ethers'


/**
 * L1 network chain IDs
 */
export enum L1ChainID {
  MAINNET = 1,
  SEPOLIA = 11155111,
  HARDHAT_LOCAL = 31337
}

/**
 * L2 network chain IDs
 */
export enum L2ChainID {
  SPECULAR = 93481,
  SPECULAR_HARDHAT_LOCAL = 31337,
}


/**
 * Enum describing the status of a message.
 */
export enum MessageStatus {
  /**
   * Message is an L1 to L2 message and has not been processed by the L2.
   */
  UNCONFIRMED_L1_TO_L2_MESSAGE,

  /**
   * Message is an L1 to L2 message and the transaction to execute the message failed.
   * When this status is returned, you will need to resend the L1 to L2 message, probably with a
   * higher gas limit.
   */
  FAILED_L1_TO_L2_MESSAGE,

  /**
   * Message is an L2 to L1 message and no state root has been published yet.
   */
  STATE_ROOT_NOT_PUBLISHED,

  /**
   * Message is ready to be proved on L1 to initiate the challenge period.
   */
  READY_TO_PROVE,

  /**
   * Message is a proved L2 to L1 message and is undergoing the challenge period.
   */
  IN_CHALLENGE_PERIOD,

  /**
   * Message is ready to be relayed.
   */
  READY_FOR_RELAY,

  /**
   * Message has been relayed.
   */
  RELAYED,
}

/**
 * Enum describing the direction of a message.
 */
export enum MessageDirection {
  L1_TO_L2,
  L2_TO_L1,
}

/**
 * Partial message that needs to be signed and executed by a specific signer.
 */
export interface CrossChainMessageRequest {
  direction: MessageDirection
  target: string
  message: string
}

/**
 * Core components of a cross chain message.
 */
export interface CoreCrossChainMessage {
  sender: string
  target: string
  message: string
  messageNonce: BigNumber
  value: BigNumber
  minGasLimit: BigNumber
}

/**
 * Describes a message that is sent between L1 and L2. Direction determines where the message was
 * sent from and where it's being sent to.
 */
export interface CrossChainMessage extends CoreCrossChainMessage {
  direction: MessageDirection
  logIndex: number
  blockNumber: number
  transactionHash: string
}

/**
 * Describes messages sent inside the L2ToL1MessagePasser on L2. Happens to be the same structure
 * as the CoreCrossChainMessage so we'll reuse the type for now.
 */
export type LowLevelMessage = CoreCrossChainMessage

/**
 * Describes a token withdrawal or deposit, along with the underlying raw cross chain message
 * behind the deposit or withdrawal.
 */
export interface TokenBridgeMessage {
  direction: MessageDirection
  from: string
  to: string
  l1Token: string
  l2Token: string
  amount: BigNumber
  data: string
  logIndex: number
  blockNumber: number
  transactionHash: string
}


/**
 * Enum describing the status of a CrossDomainMessage message receipt.
 */
export enum MessageReceiptStatus {
  RELAYED_FAILED,
  RELAYED_SUCCEEDED,
}

/**
 * CrossDomainMessage receipt.
 */
export interface MessageReceipt {
  receiptStatus: MessageReceiptStatus
  transactionReceipt: TransactionReceipt
}

/**
 * Header for a state root batch.
 */
export interface StateRootBatchHeader {
  batchIndex: BigNumber
  batchRoot: string
  batchSize: BigNumber
  prevTotalElements: BigNumber
  extraData: string
}

/**
 * Information about a state root, including header, block number, and root iself.
 */
export interface StateRoot {
  stateRoot: string
  stateRootIndexInBatch: number
  batch: StateRootBatch
}

/**
 * Information about a batch of state roots.
 */
export interface StateRootBatch {
  blockNumber: number
  header: StateRootBatchHeader
  stateRoots: string[]
}

/**
 * Proof data required to finalize an L2 to L1 message.
 */
export interface CrossChainMessageProof {
  stateRoot: string
  stateRootBatchHeader: StateRootBatchHeader
  stateRootProof: {
    index: number
    siblings: string[]
  }
  stateTrieWitness: string
  storageTrieWitness: string
}

/**
 * Stuff that can be coerced into a transaction.
 */
export type TransactionLike = string | TransactionReceipt | TransactionResponse

/**
 * Stuff that can be coerced into a CrossChainMessage.
 */
export type MessageLike =
  | CrossChainMessage
  | TransactionLike
  | TokenBridgeMessage

/**
 * Stuff that can be coerced into a CrossChainMessageRequest.
 */
export type MessageRequestLike =
  | CrossChainMessageRequest
  | CrossChainMessage
  | TransactionLike
  | TokenBridgeMessage

/**
 * Stuff that can be coerced into a provider.
 */
export type ProviderLike = string | Provider

/**
 * Stuff that can be coerced into a signer.
 */
export type SignerLike = string | Signer

/**
 * Stuff that can be coerced into a signer or provider.
 */
export type SignerOrProviderLike = SignerLike | ProviderLike

/**
 * Stuff that can be coerced into an address.
 */
export type AddressLike = string | Contract

/**
 * Stuff that can be coerced into a number.
 */
export type NumberLike = string | number | BigNumber

/**
 * L1 contract references.
 */
export interface L1Contracts {
  L1Portal: AddressLike
  L1StandardBridge: AddressLike
  L1Rollup: AddressLike
}

/**
 * L2 contract references.
 */
export interface L2Contracts {
  UUPSPlaceholder: AddressLike
  L1Oracle: AddressLike
  L2Portal: AddressLike
  L2StandardBridge: AddressLike
  L1FeeVault: AddressLike
  L2BaseFeeVault: AddressLike
}
export interface ContractsLike {
  l1: L1Contracts
  l2: L2Contracts
}

