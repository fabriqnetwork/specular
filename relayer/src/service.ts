import {
  CrossDomainMessage,
  CrossDomainMessager,
  EncodedBlockHeader,
  MessageProof,
} from "./messager";
import { ethers, BigNumber, utils } from "ethers";
import { JsonRpcProvider } from "@ethersproject/providers";

import { getStorageKey, rawBlockHeaderToEncoded } from "./utils";

import { RelayerState } from "./state";
import { DepositInitiatedEvent } from "../../contracts/typechain-types/src/bridge/IL1Portal";
import { WithdrawalInitiatedEvent } from "../../contracts/typechain-types/src/bridge/IL2Portal";
import { TxBatchAppendedEvent } from "../../contracts/typechain-types/src/ISequencerInbox";

export type RelayerConfig = {
  pollInterval: number;
  L1OracleUpdateInterval: number;

  l1ProviderEndpoint: string;
  l2ProviderEndpoint: string;
  l1RelayerPrivateKey: string;
  l2RelayerPrivateKey: string;

  inboxAddress: string;
  rollupAddress: string;
  l1OracleAddress: string;
  l1PortalAddress: string;
  l2PortalAddress: string;
};

export class RelayerService {
  readonly config: RelayerConfig;

  readonly state: RelayerState;

  readonly l1Provider: JsonRpcProvider;
  readonly l2Provider: JsonRpcProvider;
  readonly messager: CrossDomainMessager;

  constructor(config: RelayerConfig) {
    this.config = config;

    this.state = new RelayerState();

    this.l1Provider = new ethers.providers.JsonRpcProvider(
      config.l1ProviderEndpoint
    );
    this.l2Provider = new ethers.providers.JsonRpcProvider(
      config.l2ProviderEndpoint
    );

    const l1Signer = new ethers.Wallet(
      config.l1RelayerPrivateKey,
      this.l1Provider
    );
    const l2Signer = new ethers.Wallet(
      config.l2RelayerPrivateKey,
      this.l2Provider
    );
    this.messager = new CrossDomainMessager({
      l1Signer,
      l2Signer,
      inboxAddress: config.inboxAddress,
      rollupAddress: config.rollupAddress,
      l1OracleAddress: config.l1OracleAddress,
      l1PortalAddress: config.l1PortalAddress,
      l2PortalAddress: config.l2PortalAddress,
    });
  }

  async newL1BlockHeadCallback(blockNumber: number) {
    if (
      blockNumber - this.state.lastSentL1OracleBlockNumber <
      this.config.L1OracleUpdateInterval
    ) {
      return;
    }
    this.state.sentL1OracleValues(blockNumber);
    const rawBlock = await this.l1Provider.send("eth_getBlockByNumber", [
      utils.hexValue(blockNumber),
      false, // We only want the block header
    ]);
    const stateRoot = this.l1Provider.formatter.hash(rawBlock.stateRoot);
    // Sequence the state root
    await this.messager.setL1OracleValues(blockNumber, stateRoot);
    console.log("sent L1 oracle values: ", blockNumber, stateRoot);
  }

  async inboxTxBatchAppendedCallback(
    batchNumber: BigNumber,
    previousInboxSize: BigNumber,
    inboxSize: BigNumber,
    event: TxBatchAppendedEvent
  ) {
    // get l2 block number => inbox size mapping
    const tx = await event.getTransaction();
    const l2BlockNumber =
      this.messager.getLastL2BlockNumberFromAppendTxBatchCalldata(tx.data);
    this.state.updateL2BlockNumberMapping(l2BlockNumber, inboxSize);
  }

  async assertionConfirmedCallback(assertionID: BigNumber) {
    const assertion = await this.messager.getAssertion(assertionID);
    this.state.updateConfirmedInboxSize(assertionID, assertion.inboxSize);
  }

  async l1OracleValuesUpdatedCallback(blockNumber: BigNumber) {
    this.state.updatedL1OracleValues(blockNumber.toNumber());
    console.log("updated L1 oracle values: ", blockNumber.toNumber());
  }

  async depositInitiatedCallback(
    nonce: BigNumber,
    sender: string,
    target: string,
    value: BigNumber,
    gasLimit: BigNumber,
    data: string,
    depositHash: string,
    event: DepositInitiatedEvent
  ) {
    const message = {
      nonce,
      sender,
      target,
      value,
      gasLimit,
      data,
    };
    console.log("deposit initiated: ", message, event);
    this.state.addDeposit(event.blockNumber, depositHash, message);
  }

  async depositFinalizedCallback(depositHash: string, success: boolean) {
    console.log("deposit finalized: ", depositHash, success);
  }

  async withdrawalInitiatedCallback(
    nonce: BigNumber,
    sender: string,
    target: string,
    value: BigNumber,
    gasLimit: BigNumber,
    data: string,
    withdrawalHash: string,
    event: WithdrawalInitiatedEvent
  ) {
    const message = {
      nonce,
      sender,
      target,
      value,
      gasLimit,
      data,
    };
    console.log("withdrawal initiated: ", message, event);
    this.state.addWithdrawal(event.blockNumber, withdrawalHash, message);
  }

  async withdrawalFinalizedCallback(withdrawalHash: string, success: boolean) {
    console.log("withdrawal finalized: ", withdrawalHash, success);
  }

  async generateDepositProof(
    blockNumber: number,
    depositHash: string
  ): Promise<MessageProof> {
    const rawProof = await this.l1Provider.send("eth_getProof", [
      this.config.l1PortalAddress,
      [getStorageKey(depositHash)],
      utils.hexValue(blockNumber),
    ]);
    return {
      accountProof: rawProof.accountProof,
      storageProof: rawProof.storageProof.proof,
    };
  }

  async getL2EncodedBlockHeader(
    l2BlockNumber: number
  ): Promise<EncodedBlockHeader> {
    const rawBlock = await this.l2Provider.send("eth_getBlockByNumber", [
      utils.hexValue(l2BlockNumber),
      false,
    ]);
    return rawBlockHeaderToEncoded(rawBlock);
  }

  async generateWithdrawalProof(
    blockNumber: number,
    withdrawalHash: string
  ): Promise<MessageProof> {
    const rawProof = await this.l2Provider.send("eth_getProof", [
      this.config.l2PortalAddress,
      [getStorageKey(withdrawalHash)],
      utils.hexValue(blockNumber),
    ]);
    return {
      accountProof: rawProof.accountProof,
      storageProof: rawProof.storageProof.proof,
    };
  }

  async runRelayDeposits() {
    while (true) {
      // Finalize any deposits that are ready
      const nextDepositBlockNumber = this.state.getNextDepositBlockNumber();
      if (
        nextDepositBlockNumber === undefined ||
        nextDepositBlockNumber > this.state.lastUpdatedL1OracleBlockNumber
      ) {
        if (nextDepositBlockNumber !== undefined) {
          console.log(
            "waiting for L1 oracle values to update...",
            nextDepositBlockNumber,
            this.state.lastUpdatedL1OracleBlockNumber
          );
        } else {
          console.log("waiting for deposits...");
        }
        await new Promise((resolve) =>
          setTimeout(resolve, this.config.pollInterval)
        );
        continue;
      }

      const deposit = this.state.getNextDeposit()!;
      const proof = await this.generateDepositProof(
        this.state.lastUpdatedL1OracleBlockNumber,
        deposit.depositHash
      );
      try {
        await this.messager.finalizeDeposit(deposit.depositTx, proof);
      } catch (err) {
        console.error(err);
        this.state.readdDeposit(deposit);
      }
    }
  }

  async runRelayWithdrawals() {
    while (true) {
      // Finalize any withdrawals that are ready
      const nextWithdrawalBlockNumber =
        this.state.getNextWithdrawalBlockNumber();
      if (
        nextWithdrawalBlockNumber === undefined ||
        nextWithdrawalBlockNumber > this.state.lastConfirmedL2BlockNumber
      ) {
        if (nextWithdrawalBlockNumber !== undefined) {
          console.log(
            "waiting for L2 block to be confirmed...",
            nextWithdrawalBlockNumber,
            this.state.lastConfirmedL2BlockNumber
          );
        } else {
          console.log("waiting for withdrawals...");
        }
        await new Promise((resolve) =>
          setTimeout(resolve, this.config.pollInterval)
        );
        continue;
      }

      const withdrawal = this.state.getNextWithdrawal()!;
      const header = await this.getL2EncodedBlockHeader(
        this.state.lastConfirmedL2BlockNumber
      );
      const proof = await this.generateWithdrawalProof(
        this.state.lastConfirmedL2BlockNumber,
        withdrawal.withdrawalHash
      );
      try {
        await this.messager.finalizeWithdrawal(
          withdrawal.withdrawalTx,
          this.state.lastConfirmedAssertionID,
          header,
          proof
        );
      } catch (err) {
        console.error(err);
        this.state.readdWithdrawal(withdrawal);
      }
    }
  }

  async start() {
    this.l1Provider.on("block", this.newL1BlockHeadCallback.bind(this));
    this.messager.onInboxTxBatchAppend(
      this.inboxTxBatchAppendedCallback.bind(this)
    );
    this.messager.onAssertionConfirmed(
      this.assertionConfirmedCallback.bind(this)
    );
    this.messager.onL1OracleValuesUpdated(
      this.l1OracleValuesUpdatedCallback.bind(this)
    );
    this.messager.onDepositInitiated(this.depositInitiatedCallback.bind(this));
    this.messager.onDepositFinalized(this.depositFinalizedCallback.bind(this));
    this.messager.onWithdrawalInitiated(
      this.withdrawalInitiatedCallback.bind(this)
    );
    this.messager.onWithdrawalFinalized(
      this.withdrawalFinalizedCallback.bind(this)
    );

    console.log("started relayer");

    this.runRelayDeposits();
    this.runRelayWithdrawals();
  }
}
