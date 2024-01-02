import { BigNumber } from "ethers";

export type OnboardingServiceConfig = {
  l1ProviderEndpoint: string;
  l2ProviderEndpoint: string;
  l2FunderPrivateKey: string;

  l1PortalAddress: string;

  depositFundingThreshold: BigNumber;
};
