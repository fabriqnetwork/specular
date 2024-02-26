import "dotenv/config";
import { BigNumber, ethers } from "ethers";
import fastify from "fastify";

import {
  OnboardingServiceConfig,
  OnboardingService,
  relayerPlugin,
} from "./onboarding";

export default async function main() {
  const config: OnboardingServiceConfig = {
    l1ProviderEndpoint: process.env.L1_PROVIDER_ENDPOINT!,
    l2ProviderEndpoint: process.env.L2_PROVIDER_ENDPOINT!,
    l2FunderPrivateKey: process.env.L2_FUNDER_PRIVATE_KEY!,
    l1PortalAddress: process.env.L1_PORTAL_ADDRESS!,
    depositFundingThreshold: process.env.DEPOSIT_FUNDING_THRESHOLD
      ? BigNumber.from(process.env.DEPOSIT_FUNDING_THRESHOLD)
      : ethers.utils.parseEther("0.01"),
  };

  const relayer = new OnboardingService(config);
  await relayer.start();

  const app = fastify({
    logger: true,
    trustProxy: true,
  });
  app.register(relayerPlugin, { onboarding: relayer });

  await app.listen({ port: Number(process.env.APP_PORT ?? 3092) });
}
