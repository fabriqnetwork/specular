import { Signer, BigNumberish, BigNumber } from "ethers";
import { NonceManager } from "@ethersproject/experimental";

import { TypedListener } from "../../../contracts/typechain-types/common";
import {
  L1Oracle,
  L1Oracle__factory,
  IL2Portal,
  IL2Portal__factory,
} from "../../../contracts/typechain-types";
import { L1OracleValuesUpdatedEvent } from "../../../contracts/typechain-types/src/bridge/L1Oracle";
import { CrossDomainMessage } from ".";

export type CrossDomainMessagerConfig = {
  l2Relayer: Signer;
  l2Funder: Signer;

  l1OracleAddress: string;
  l2PortalAddress: string;
};

export type MessageProof = {
  accountProof: string[];
  storageProof: string[];
};

export class CrossDomainMessager {
  readonly config: CrossDomainMessagerConfig;

  readonly l2Relayer: NonceManager;
  readonly l2Funder: NonceManager;

  readonly l1Oracle: L1Oracle;
  readonly l2Portal: IL2Portal;

  static async create(config: CrossDomainMessagerConfig) {
    const l2Relayer = new NonceManager(config.l2Relayer);
    let l2Funder: NonceManager;
    if (
      (await config.l2Funder.getAddress()) !==
      (await config.l2Relayer.getAddress())
    ) {
      l2Funder = new NonceManager(config.l2Funder);
    } else {
      l2Funder = l2Relayer;
    }
    return new CrossDomainMessager(
      config,
      new NonceManager(l2Relayer),
      new NonceManager(l2Funder)
    );
  }

  constructor(
    config: CrossDomainMessagerConfig,
    l2Relayer: NonceManager,
    l2Funder: NonceManager
  ) {
    this.config = config;

    this.l2Relayer = l2Relayer;
    this.l2Funder = l2Funder;

    this.l1Oracle = L1Oracle__factory.connect(
      config.l1OracleAddress,
      this.l2Relayer
    );
    this.l2Portal = IL2Portal__factory.connect(
      config.l2PortalAddress,
      this.l2Funder
    );
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
    return await tx.wait();
  }

  onL1OracleValuesUpdated(callback: TypedListener<L1OracleValuesUpdatedEvent>) {
    this.l1Oracle.on(this.l1Oracle.filters.L1OracleValuesUpdated(), callback);
  }
}
