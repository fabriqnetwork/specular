import fs from "fs";
import { ethers } from "ethers";
import hre from "hardhat";
import { parseFlag } from "./utils";

type RawLog = {
  topics: string[],
  data: string
}

async function main() {
  const baseConfigPath = parseFlag("--in");
  const configPath = parseFlag("--out");
  const genesisPath = parseFlag("--genesis");
  const genesisHashPath = parseFlag("--genesis-hash-path");
  const deploymentsPath = parseFlag("--deployments", "./deployments/localhost");
  await generateConfigFile(baseConfigPath, configPath, genesisPath, genesisHashPath, deploymentsPath);
}

/**
 * Reads the L1 and L2 genesis block info from the specified deployment and
 * adds it to the base config file
 */
export async function generateConfigFile(
  baseConfigPath: string,
  configPath: string,
  genesisPath: string,
  genesisHashPath: string,
  deploymentsPath: string
) {
  // check the deployments dir - error out if it is not there
  const contract = "Proxy__Rollup";
  const deployment = JSON.parse(fs.readFileSync(`${deploymentsPath}/${contract}.json`, "utf-8"))

  // exctract L1 block hash and L1 block number from receipt
  const l1Number = deployment.receipt.blockNumber;
  const l1Hash = deployment.receipt.blockHash;

  // Parse genesis hash file to get L2 genesis hash
  const l2Hash = JSON.parse(fs.readFileSync(genesisHashPath, "utf-8")).hash;

  // Write out new file
  // TODO: use on-chain data-only or genesis-only
  const baseConfig = JSON.parse(fs.readFileSync(baseConfigPath, "utf-8"));
  baseConfig.genesis.l1.hash = l1Hash;
  baseConfig.genesis.l1.number = l1Number;
  baseConfig.genesis.l2.hash = l2Hash;
  const genesis = JSON.parse(fs.readFileSync(genesisPath, "utf-8"));
  baseConfig.genesis.l2_time = ethers.BigNumber.from(genesis.timestamp).toNumber() || 0;

  fs.writeFileSync(configPath, JSON.stringify(baseConfig, null, 2));
  console.log(`successfully wrote config to: ${configPath}`)
}

if (!require.main!.loaded) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
