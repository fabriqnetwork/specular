import { BigNumber, Signer } from "ethers";
import { NonceManager } from "@ethersproject/experimental";

import {
  IL2Portal,
  IL2Portal__factory,
  L1Oracle,
  L1Oracle__factory,
  l1OracleAddress,
  l2PortalAddress,
} from "@specularl2/sdk";
import { CrossDomainMessage } from ".";

export type CrossDomainMessagerConfig = {
  l2Funder: Signer;
};

export type MessageProof = {
  accountProof: string[];
  storageProof: string[];
};

export class CrossDomainMessenger {
  readonly config: CrossDomainMessagerConfig;

  readonly l2Funder: NonceManager;

  readonly l1Oracle: L1Oracle;
  readonly l2Portal: IL2Portal;

  static async create(config: CrossDomainMessagerConfig) {
    const l2Funder = new NonceManager(config.l2Funder);
    return new CrossDomainMessenger(config, new NonceManager(l2Funder));
  }

  constructor(config: CrossDomainMessagerConfig, l2Funder: NonceManager) {
    this.config = config;

    this.l2Funder = l2Funder;

    this.l1Oracle = L1Oracle__factory.connect(l1OracleAddress, this.l2Funder);
    this.l2Portal = IL2Portal__factory.connect(l2PortalAddress, this.l2Funder);
  }

  async getL1OracleBlockNumber() {
    return await this.l1Oracle.number();
  }

  async finalizeDeposit(
    l1BlockNumber: BigNumber,
    depositTx: CrossDomainMessage,
    depositProof: MessageProof
  ) {
    const tx = await this.l2Portal.finalizeDepositTransaction(
      l1BlockNumber,
      depositTx,
      depositProof.accountProof,
      depositProof.storageProof
    );
    return await tx.wait();
  }
}
