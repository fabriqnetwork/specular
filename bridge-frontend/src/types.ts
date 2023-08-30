import type { ethers,BigNumber } from "ethers";

export type Data = {
  status: string;
  error?: string;
  dataFrom?: ethers.providers.TransactionResponse;
  dataToHash?: string;
  l1BlockNumber?: number;
  proofL1BlockNumber?:number;
  crossDomainMessage?:CrossDomainMessage;
}

export type CrossDomainMessage = {
  version: number;
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

export type wallet = {
  address: string;
  chainId: number;
  provider: any;
}
