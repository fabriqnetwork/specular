import {
  Provider,
  TransactionReceipt,
  TransactionResponse,
} from "@ethersproject/abstract-provider";
import { Signer } from "@ethersproject/abstract-signer";
import { Contract, BigNumber } from "ethers";

/**
 * L1 network chain IDs.
 */
export enum L1ChainID {
  MAINNET = 1,
  SEPOLIA = 11155111,
  HARDHAT_LOCAL = 31337,
}

/**
 * Message finalizaton status.
 */
export enum MessageStatus {
  // not ready to be finalized
  PENDING = 0,
  // ready to be finalized
  READY = 1,
  // already finalized
  DONE = 2,
}

/**
 * L2 network chain IDs.
 */
export enum L2ChainID {
  SPECULAR = 93481,
  SPECULAR_HARDHAT_LOCAL = 13527,
}

/**
 * Stuff that can be coerced into a transaction.
 */
export type TransactionLike = string | TransactionReceipt | TransactionResponse;

/**
 * Stuff that can be coerced into a provider.
 */
export type ProviderLike = string | Provider;

/**
 * Stuff that can be coerced into a signer.
 */
export type SignerLike = string | Signer;

/**
 * Stuff that can be coerced into a signer or provider.
 */
export type SignerOrProviderLike = SignerLike | ProviderLike;

/**
 * Stuff that can be coerced into an address.
 */
export type AddressLike = string | Contract;

/**
 * Stuff that can be coerced into a number.
 */
export type NumberLike = string | number | BigNumber;

/**
 * L1 contract references.
 */
export interface L1Contracts {
  L1Portal: AddressLike;
  L1StandardBridge: AddressLike;
  L1Rollup: AddressLike;
}

/**
 * L2 contract references.
 */
export interface L2Contracts {
  UUPSPlaceholder: AddressLike;
  L1Oracle: AddressLike;
  L2Portal: AddressLike;
  L2StandardBridge: AddressLike;
  L1FeeVault: AddressLike;
  L2BaseFeeVault: AddressLike;
}

/**
 * L1 and L2 contracts references.
 */
export interface ContractsLike {
  l1: L1Contracts;
  l2: L2Contracts;
}

export enum MessageType {
  DEPOSIT = "Deposit",
  WITHDRAWAL = "Withdrawal",
}

export type Message = {
  version: bigint;
  nonce: bigint;
  sender: AddressLike;
  target: AddressLike;
  value: bigint;
  gasLimit: bigint;
  data: string;
};

export type BridgeTransaction = {
  messageHash: string;
  block: NumberLike;
  amount: bigint;
  type: MessageType;
  action: {
    status: MessageStatus;
    message: Message;
    contract: any;
    chain: any;
  };
};
