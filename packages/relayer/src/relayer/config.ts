import { BigNumber } from "ethers";

export type RelayerConfig = {
  L1OracleUpdateInterval: number;

  l1ProviderEndpoint: string;
  l2ProviderEndpoint: string;
  l2RelayerPrivateKey: string;
  l2FunderPrivateKey: string;

  l1OracleAddress: string;
  l1PortalAddress: string;
  l2PortalAddress: string;

  depositFundingThreshold: BigNumber;
};
