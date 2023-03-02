import "dotenv/config";
import express from "express";
import { RelayerService } from "./src/service";

async function main() {
  const relayerConfig = {
    pollInterval: Number(process.env.POLL_INTERVAL) ?? 1000,
    L1OracleUpdateInterval:
      Number(process.env.L1_ORACLE_UPDATE_INTERVAL) ?? 1000,

    l1ProviderEndpoint: process.env.L1_PROVIDER_ENDPOINT!,
    l2ProviderEndpoint: process.env.L2_PROVIDER_ENDPOINT!,
    l1RelayerPrivateKey: process.env.L1_RELAYER_PRIVATE_KEY!,
    l2RelayerPrivateKey: process.env.L2_RELAYER_PRIVATE_KEY!,

    inboxAddress: process.env.INBOX_ADDRESS!,
    rollupAddress: process.env.ROLLUP_ADDRESS!,
    l1OracleAddress: process.env.L1_ORACLE_ADDRESS!,
    l1PortalAddress: process.env.L1_PORTAL_ADDRESS!,
    l2PortalAddress: process.env.L2_PORTAL_ADDRESS!,
  };
  const relayer = new RelayerService(relayerConfig);
  await relayer.start();

  const app = express();
  app.get("/", (req, res) => {
    res.send("Hello World!");
  });
  app.listen(process.env.APP_PORT);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
