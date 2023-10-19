import "dotenv/config";
import { BigNumber, ethers } from "ethers";
import fastify from "fastify";

import { RelayerConfig, RelayerService, relayerPlugin } from "./relayer";

export default async function main() {
  const relayerConfig: RelayerConfig = {
    L1OracleUpdateInterval: Number(process.env.L1_ORACLE_UPDATE_INTERVAL) ?? 10,

    l1ProviderEndpoint: process.env.L1_PROVIDER_ENDPOINT!,
    l2ProviderEndpoint: process.env.L2_PROVIDER_ENDPOINT!,
    l2RelayerPrivateKey: process.env.L2_RELAYER_PRIVATE_KEY!,
    l2FunderPrivateKey:
      process.env.L2_FUNDER_PRIVATE_KEY ?? process.env.L2_RELAYER_PRIVATE_KEY!,

    l1OracleAddress: process.env.L1_ORACLE_ADDRESS!,
    l1PortalAddress: process.env.L1_PORTAL_ADDRESS!,
    l2PortalAddress: process.env.L2_PORTAL_ADDRESS!,

    depositFundingThreshold:
      process.env.DEPOSIT_FUNDING_THRESHOLD ? BigNumber.from(process.env.DEPOSIT_FUNDING_THRESHOLD) :
      ethers.utils.parseEther("0.01"),
  };

  const relayer = new RelayerService(relayerConfig);
  await relayer.start();

  const app = fastify({
    logger: true,
    trustProxy: true,
  });
  app.register(relayerPlugin, { relayer });

  await app.listen({ port: Number(process.env.APP_PORT ?? 3092) });
}
