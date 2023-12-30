import path from "node:path";
import fs from "fs";
import { parseFlag } from "./utils";
import hre from "hardhat";

const { ethers } = hre;

const CONTRACTS_DIR = path.join(__dirname, "/../../")
require("dotenv").config({ path: path.join(CONTRACTS_DIR, ".genesis.env")});
const EXPORTED_PATH = process.env.GENESIS_EXPORTED_HASH_PATH || "";

async function main() {
  const exported = require(path.join(CONTRACTS_DIR, EXPORTED_PATH));
  console.log({ exported })
  const initialBlockHash = (exported.hash || "") as string;
  if (!initialBlockHash) {
     throw Error(`blockHash not found\n$`);
  }
  console.log("initial blockHash:", initialBlockHash);

  const initialStateRoot = (exported.stateRoot || "") as string;
  if (!initialStateRoot) {
     throw Error(`stateRoot not found\n$`);
  }
  console.log("initial stateRoot:", initialStateRoot);

  const initialRollupState = {
    assertionID: 0,
    l2BlockNum: 0,
    l2BlockHash: initialBlockHash,
    l2StateRoot: initialStateRoot,
  };

  // check the deployments dir - error out if it is not there
  const deploymentsPath = parseFlag("--deployments", "./deployments/localhost");
  const deployment = JSON.parse(
    fs.readFileSync(`${deploymentsPath}/Proxy__Rollup.json`, "utf-8"),
  );

  const { deployer } = await hre.getNamedAccounts();
  const RollupFactory = await ethers.getContractFactory("Rollup", deployer);
  const rollup = RollupFactory.attach(deployment.address);
  const tx = await rollup.initializeGenesis(initialRollupState);
  await tx.wait();
}

if (!require.main!.loaded) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
