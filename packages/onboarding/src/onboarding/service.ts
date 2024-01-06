import { ethers, BigNumber } from "ethers";
import { JsonRpcProvider } from "@ethersproject/providers";

import {
  OnboardingServiceConfig as OnboardingConfig,
  CrossDomainMessenger,
  CrossDomainMessage,
  MessageProof,
  getStorageKey,
} from ".";

export class OnboardingService {
  readonly config: OnboardingConfig;

  readonly l1Provider: JsonRpcProvider;
  readonly l2Provider: JsonRpcProvider;

  messenger: CrossDomainMessenger | undefined = undefined;

  readonly l2Funder: ethers.Wallet;

  constructor(config: OnboardingConfig) {
    this.config = config;

    this.l1Provider = new ethers.providers.JsonRpcProvider(
      config.l1ProviderEndpoint
    );
    this.l2Provider = new ethers.providers.JsonRpcProvider(
      config.l2ProviderEndpoint
    );

    this.l2Funder = new ethers.Wallet(
      config.l2FunderPrivateKey,
      this.l2Provider
    );
  }

  async start() {
    this.messenger = await CrossDomainMessenger.create({
      l2Funder: this.l2Funder,
    });
    console.log("started onboarding service");
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

  async fundDeposit(depositTx: CrossDomainMessage, depositHash: string) {
    if (this.messenger === undefined) {
      throw new Error("Messenger not started");
    }
    const blockNumber = await this.messenger.getL1OracleBlockNumber();
    const proof = await this.generateDepositProof(depositHash, blockNumber);
    console.log(proof);
    const tx = await this.messenger.finalizeDeposit(
      blockNumber,
      depositTx,
      proof
    );
    return tx.transactionHash;
  }
}
