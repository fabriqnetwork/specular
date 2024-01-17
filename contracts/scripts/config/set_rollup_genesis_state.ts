import path from "node:path";
import fs from "fs";
import { parseFlag } from "./utils";
import hre from "hardhat";

const { ethers } = hre;

const CONTRACTS_DIR = path.join(__dirname, "/../../");
require("dotenv").config({ path: path.join(CONTRACTS_DIR, ".genesis.env") });
const EXPORTED_PATH = process.env.GENESIS_EXPORTED_HASH_PATH || "";
const ROLLUP_ADDR = process.env.ROLLUP_ADDR || "invalid address";

async function main() {
  const exported = require(path.join(CONTRACTS_DIR, EXPORTED_PATH));
  console.log({ exported });
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

  const { deployer } = await hre.getNamedAccounts();
  const RollupFactory = await ethers.getContractFactory("Rollup", deployer);
  const rollup = RollupFactory.attach(ROLLUP_ADDR);
  const tx = await rollup.initializeGenesis(initialRollupState);
  await tx.wait();
  console.log("initialized genesis state on rollup contract");
}

if (!require.main!.loaded) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
