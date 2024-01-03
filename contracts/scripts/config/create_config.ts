import fs from "fs";
import { ethers } from "ethers";
import { parseFlag } from "./utils";

require("dotenv").config();

// TODO: consider moving to golang (ops).
async function main() {
  const baseConfigPath = parseFlag("--in");
  const configPath = parseFlag("--out");
  const genesisPath = parseFlag("--genesis");
  const genesisHashPath = parseFlag("--genesis-hash-path");
  const deploymentsPath = parseFlag("--deployments", "./deployments/localhost");
  await generateConfigFile(
    configPath,
    baseConfigPath,
    genesisPath,
    genesisHashPath,
    deploymentsPath,
  );
}

/**
 * Reads the L1 and L2 genesis block info from the specified deployment and
 * adds it to the base config file, along with other config params.
 * Outputs the new config file at `configPath`.
 */
export async function generateConfigFile(
  configPath: string,
  baseConfigPath: string,
  genesisPath: string,
  genesisHashPath: string,
  deploymentsPath: string,
) {
  // check the deployments dir - error out if it is not there
  const contract = "Proxy__Rollup";
  const deployment = JSON.parse(
    fs.readFileSync(`${deploymentsPath}/${contract}.json`, "utf-8"),
  );
  const inboxDeployment = JSON.parse(
    fs.readFileSync(`${deploymentsPath}/Proxy__SequencerInbox.json`, "utf-8"),
  );

  // extract L1 block hash and L1 block number from receipt
  const l1Number = deployment.receipt.blockNumber;
  const l1Hash = deployment.receipt.blockHash;

  // Parse genesis and hash file.
  const l2Hash = JSON.parse(fs.readFileSync(genesisHashPath, "utf-8")).hash;
  const genesis = JSON.parse(fs.readFileSync(genesisPath, "utf-8"));
  // Set genesis fields.
  const baseConfig = JSON.parse(fs.readFileSync(baseConfigPath, "utf-8"));
  baseConfig.genesis = {
    l1: {
      hash: l1Hash,
      number: l1Number,
    },
    l2: {
      hash: l2Hash,
      number: ethers.BigNumber.from(genesis.number).toNumber(),
    },
    l2_time: ethers.BigNumber.from(genesis.timestamp).toNumber(),
    system_config: {
      batcherAddr: process.env.SEQUENCER_ADDRESS,
      gasLimit: ethers.BigNumber.from(genesis.gasLimit).toNumber(),
      overhead:
        "0x0000000000000000000000000000000000000000000000000000000000000000",
      scalar:
        "0x0000000000000000000000000000000000000000000000000000000000000000",
    },
  };
  // Set other fields.
  baseConfig.l2_chain_id = genesis.config.chainId;
  baseConfig.batch_inbox_address = inboxDeployment.address;
  baseConfig.rollup_address = deployment.address;

  // Write out new file.
  fs.writeFileSync(configPath, JSON.stringify(baseConfig, null, 2));
  console.log(`successfully wrote config to: ${configPath}`);
}

if (!require.main!.loaded) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
