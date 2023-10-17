import { task } from "hardhat/config";
import { getProxyName } from "../deploy/utils";
import { readFileSync, write, writeFileSync } from "fs";

// Usage: npx hardhat --network <l1-network> inject-l1-config --file <path-to-rollup.json>
task("inject-l1-config", "Inject L1 deployment info to Rollup.json config file")
  .addParam("file", "Path to Rollup.json config file")
  .setAction(async (taskArgs, hre) => {
    const configFile = JSON.parse(readFileSync(taskArgs.file, "utf8"));

    const deployment = await hre.deployments.get(
      getProxyName("SequencerInbox")
    );
    const l1BlockHash = deployment.receipt?.blockHash;
    if (!l1BlockHash) {
      throw Error("Could not find L1 block hash of SequencerInbox deployment");
    }
    const block = await hre.ethers.provider.getBlock(l1BlockHash);

    if (!configFile["Genesis"]) {
      configFile["Genesis"] = {};
    }
    if (!configFile["Genesis"]["L1"]) {
      configFile["Genesis"]["L1"] = {};
    }
    configFile["Genesis"]["L1"]["Hash"] = l1BlockHash;
    configFile["Genesis"]["L1"]["Number"] = block.number;
    configFile["Genesis"]["L2Time"] = block.timestamp;
    writeFileSync(taskArgs.file, JSON.stringify(configFile));
  });
