import type { BigNumber } from "ethers";

export type CrossDomainMessage = {
  nonce: BigNumber;
  sender: string;
  target: string;
  value: BigNumber;
  gasLimit: BigNumber;
  data: string;
};

export type PendingDeposit = {
  l1BlockNumber: number;
  proofL1BlockNumber: number | undefined;
  depositHash: string;
  depositTx: CrossDomainMessage;
};

export type PendingWithdrawal = {
  l2BlockNumber: number;
  proofL2BlockNumber: number | undefined;
  inboxSize: BigNumber | undefined;
  assertionID: BigNumber | undefined;
  withdrawalHash: string;
  withdrawalTx: CrossDomainMessage;
};

export type MessageProof = {
  accountProof: string[];
  storageProof: string[];
};
