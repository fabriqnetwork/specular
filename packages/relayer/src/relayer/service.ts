import { ethers, BigNumber, utils } from "ethers";
import { JsonRpcProvider, WebSocketProvider } from "@ethersproject/providers";

import {
  RelayerConfig,
  RelayerState,
  CrossDomainMessager,
  CrossDomainMessage,
  MessageProof,
  getStorageKey,
} from ".";

export class RelayerService {
  readonly config: RelayerConfig;

  readonly state: RelayerState;

  readonly l1Provider: JsonRpcProvider;
  readonly l2Provider: JsonRpcProvider;

  messager: CrossDomainMessager | undefined = undefined;

  readonly l2Relayer: ethers.Wallet;
  readonly l2Funder: ethers.Wallet;

  constructor(config: RelayerConfig) {
    this.config = config;
    this.state = new RelayerState();

    this.l1Provider = new ethers.providers.JsonRpcProvider(
      config.l1ProviderEndpoint
    );
    this.l2Provider = new ethers.providers.JsonRpcProvider(
      config.l2ProviderEndpoint
    );

    this.l2Relayer = new ethers.Wallet(
      config.l2RelayerPrivateKey,
      this.l2Provider
    );
    this.l2Funder = new ethers.Wallet(
      config.l2FunderPrivateKey,
      this.l2Provider
    );
  }

  async start() {
    this.messager = await CrossDomainMessager.create({
      l2Relayer: this.l2Relayer,
      l2Funder: this.l2Funder,
      l1OracleAddress: this.config.l1OracleAddress,
      l2PortalAddress: this.config.l2PortalAddress,
    });
    this.l1Provider.on("block", this.newL1BlockHeadCallback.bind(this));
    this.messager.onL1OracleValuesUpdated(
      this.l1OracleValuesUpdatedCallback.bind(this)
    );
    console.log("started relayer");
  }

  async newL1BlockHeadCallback(blockNumber: number | BigNumber) {
    if (this.messager === undefined) {
      return;
    }
    blockNumber = BigNumber.from(blockNumber);
    if (
      blockNumber
        .sub(this.state.lastSentL1OracleBlockNumber)
        .lt(this.config.L1OracleUpdateInterval)
    ) {
      return;
    }
    const oldValue = this.state.lastSentL1OracleBlockNumber;
    try {
      this.state.sentL1OracleValues(blockNumber);
      const rawBlock = await this.l1Provider.send("eth_getBlockByNumber", [
        utils.hexValue(blockNumber),
        false, // We only want the block header
      ]);
      const stateRoot = this.l1Provider.formatter.hash(rawBlock.stateRoot);
      // Sequence the state root
      await this.messager.setL1OracleValues(blockNumber, stateRoot);
      console.log("sent L1 oracle values: ", blockNumber, stateRoot);
    } catch (e) {
      console.log("failed to send L1 oracle values: ", e);
      this.state.lastSentL1OracleBlockNumber = oldValue;
    }
  }

  async l1OracleValuesUpdatedCallback(blockNumber: BigNumber) {
    this.state.updatedL1OracleValues(blockNumber);
    console.log("updated L1 oracle values: ", blockNumber.toString());
  }

  async generateDepositProof(
    depositHash: string,
    blockNumber: BigNumber
  ): Promise<MessageProof> {
    const rawProof = await this.l1Provider.send("eth_getProof", [
      this.config.l1PortalAddress,
      [getStorageKey(depositHash)],
      ethers.utils.hexValue(blockNumber),
    ]);
    return {
      accountProof: rawProof.accountProof,
      storageProof: rawProof.storageProof[0].proof,
    };
  }

  async fundDeposit(
    depositTx: CrossDomainMessage,
    depositHash: string,
    blockNumber: BigNumber
  ) {
    if (this.messager === undefined) {
      throw new Error("Relayer not started");
    }
    const proof = await this.generateDepositProof(depositHash, blockNumber);
    console.log(proof);
    const tx = await this.messager.finalizeDeposit(depositTx, proof);
    return tx.transactionHash;
  }
}
