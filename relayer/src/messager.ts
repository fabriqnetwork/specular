import {
  L1Oracle,
  L1Oracle__factory,
  IL1Portal,
  IL1Portal__factory,
  IL2Portal,
  IL2Portal__factory,
  ISequencerInbox,
  ISequencerInbox__factory,
  IRollup,
  IRollup__factory,
} from "../../contracts/typechain-types";

import { Signer, BigNumberish, BigNumber } from "ethers";
import { NonceManager } from "@ethersproject/experimental";

import { L1OracleValuesUpdatedEvent } from "../../contracts/typechain-types/src/bridge/L1Oracle";
import {
  WithdrawalFinalizedEvent,
  DepositInitiatedEvent,
} from "../../contracts/typechain-types/src/bridge/IL1Portal";
import {
  DepositFinalizedEvent,
  WithdrawalInitiatedEvent,
} from "../../contracts/typechain-types/src/bridge/IL2Portal";

import { TypedListener } from "../../contracts/typechain-types/common";
import { TxBatchAppendedEvent } from "../../contracts/typechain-types/src/ISequencerInbox";
import {
  AssertionConfirmedEvent,
  AssertionCreatedEvent,
} from "../../contracts/typechain-types/src/IRollup";

export type CrossDomainMessage = {
  nonce: BigNumber;
  sender: string;
  target: string;
  value: BigNumber;
  gasLimit: BigNumber;
  data: string;
};

export type CrossDomainMessagerConfig = {
  l1Signer: Signer;
  l2Signer: Signer;

  inboxAddress: string;
  rollupAddress: string;
  l1OracleAddress: string;
  l1PortalAddress: string;
  l2PortalAddress: string;
};

export type EncodedBlockHeader = string;

export type MessageProof = {
  accountProof: string[];
  storageProof: string[];
};

export class CrossDomainMessager {
  readonly config: CrossDomainMessagerConfig;

  readonly l1Signer: NonceManager;
  readonly l2Signer: NonceManager;

  readonly inbox: ISequencerInbox;
  readonly rollup: IRollup;
  readonly l1Oracle: L1Oracle;
  readonly l1Portal: IL1Portal;
  readonly l2Portal: IL2Portal;

  constructor(config: CrossDomainMessagerConfig) {
    this.config = config;

    this.l1Signer = new NonceManager(config.l1Signer);
    this.l2Signer = new NonceManager(config.l2Signer);

    this.inbox = ISequencerInbox__factory.connect(
      config.inboxAddress,
      this.l1Signer
    );
    this.rollup = IRollup__factory.connect(config.rollupAddress, this.l1Signer);
    this.l1Oracle = L1Oracle__factory.connect(
      config.l1OracleAddress,
      this.l2Signer
    );
    this.l1Portal = IL1Portal__factory.connect(
      config.l1PortalAddress,
      this.l1Signer
    );
    this.l2Portal = IL2Portal__factory.connect(
      config.l2PortalAddress,
      this.l2Signer
    );
  }

  getLastL2BlockNumberFromAppendTxBatchCalldata(data: string): number {
    const decoded = this.inbox.interface.decodeFunctionData(
      this.inbox.interface.functions[
        "appendTxBatch(uint256[],uint256[],bytes)"
      ],
      data
    );
    const contexts: BigNumber[] = decoded[0];
    const lastL2BlockNumber = contexts[contexts.length - 2].toNumber();
    return lastL2BlockNumber;
  }

  async getAssertion(assertionID: BigNumber) {
    return await this.rollup.getAssertion(assertionID);
  }

  async setL1OracleValues(blockNumber: BigNumberish, stateRoot: string) {
    const tx = await this.l1Oracle.setL1OracleValues(blockNumber, stateRoot);
    await tx.wait();
  }

  async finalizeDeposit(
    depositTx: CrossDomainMessage,
    depositProof: MessageProof
  ) {
    const tx = await this.l2Portal.finalizeDepositTransaction(
      depositTx,
      depositProof.accountProof,
      depositProof.storageProof
    );
    await tx.wait();
  }

  async finalizeWithdrawal(
    withdrawalTx: CrossDomainMessage,
    assertionID: BigNumberish,
    l2GasUsed: BigNumberish,
    vmHash: string,
    // encodedBlockHeader: string,
    withdrawalProof: MessageProof
  ) {
    const tx = await this.l1Portal.finalizeWithdrawalTransaction(
      withdrawalTx,
      assertionID,
      l2GasUsed,
      vmHash,
      // encodedBlockHeader,
      withdrawalProof.accountProof,
      withdrawalProof.storageProof
    );
    await tx.wait();
  }

  onInboxTxBatchAppend(callback: TypedListener<TxBatchAppendedEvent>) {
    this.inbox.on(this.inbox.filters.TxBatchAppended(), callback);
  }

  onAssertionCreated(callback: TypedListener<AssertionCreatedEvent>) {
    this.rollup.on(this.rollup.filters.AssertionCreated(), callback);
  }

  onAssertionConfirmed(callback: TypedListener<AssertionConfirmedEvent>) {
    this.rollup.on(this.rollup.filters.AssertionConfirmed(), callback);
  }

  onL1OracleValuesUpdated(callback: TypedListener<L1OracleValuesUpdatedEvent>) {
    this.l1Oracle.on(this.l1Oracle.filters.L1OracleValuesUpdated(), callback);
  }

  onDepositInitiated(callback: TypedListener<DepositInitiatedEvent>) {
    this.l1Portal.on(this.l1Portal.filters.DepositInitiated(), callback);
  }

  onWithdrawalInitiated(callback: TypedListener<WithdrawalInitiatedEvent>) {
    this.l2Portal.on(this.l2Portal.filters.WithdrawalInitiated(), callback);
    this.l2Portal.interface.events["DepositFinalized(bytes32,bool)"];
  }

  onDepositFinalized(callback: TypedListener<DepositFinalizedEvent>) {
    this.l2Portal.on(this.l2Portal.filters.DepositFinalized(), callback);
  }

  onWithdrawalFinalized(callback: TypedListener<WithdrawalFinalizedEvent>) {
    this.l1Portal.on(this.l1Portal.filters.WithdrawalFinalized(), callback);
  }
}
