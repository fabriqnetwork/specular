import path from "node:path";
import fs from "fs";
import { parseFlag } from "./utils";
import hre from "hardhat";

const { ethers } = hre;

const CONTRACTS_DIR = path.join(__dirname, "/../../");
require("dotenv").config({ path: path.join(CONTRACTS_DIR, ".genesis.env") });
const EXPORTED_PATH = process.env.GENESIS_EXPORTED_HASH_PATH || "";
const ROLLUP_ADDR = process.env.ROLLUP_ADDR || "invalid address";
const SEQUENCER_ADDR = process.env.SEQUENCER_INBOX_ADDR || "invalid address";
const VERIFIER_ADDR = process.env.VERIFIER_ADDR || "invalid address";
const VAULT_ADDR = process.env.DEPLOYER_ADDRESS || "invalid address";
const VALIDATOR_ADDR = process.env.VALIDATOR_ADDRESS || "invalid address";

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
  var tx = await rollup.initializeGenesis(initialRollupState);
  await tx.wait();
  console.log("initialized genesis state on rollup contract");

  console.log("configuring rollup contract...");
  const confirmationPeriod = 12; // TODO: move to config
  const challengePeriod = 0;
  const minimumAssertionPeriod = 0;
  const baseStakeAmount = 0;

  tx = await rollup.setDAProvider(SEQUENCER_ADDR);
  await tx.wait();
  console.log("set DA Provider");

  tx = await rollup.setVerifier(VERIFIER_ADDR);
  await tx.wait();
  console.log("set verifier");

  tx = await rollup.setVault(VAULT_ADDR);
  await tx.wait();
  console.log("set vault");

  tx = await rollup.setConfirmationPeriod(confirmationPeriod);
  await tx.wait();
  console.log("set confirmation period");

  tx = await rollup.setChallengePeriod(challengePeriod);
  await tx.wait();
  console.log("set challenge period");

  tx = await rollup.setMinimumAssertionPeriod(minimumAssertionPeriod);
  await tx.wait();
  console.log("set minimum assertion period");

  tx = await rollup.setBaseStakeAmount(baseStakeAmount);
  await tx.wait();
  console.log("set base stake amount");

  tx = await rollup.addValidator(VALIDATOR_ADDR);
  await tx.wait();
  console.log("set validator");

  console.log("configured rollup contract");

  tx = await rollup.unpause();
  await tx.wait();

  console.log("rollup contract ready");

}

if (!require.main!.loaded) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
