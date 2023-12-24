import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployer } = await hre.getNamedAccounts();
  const args = [process.env.SEQUENCER_ADDRESS];
  console.log("Deploying SequencerInbox with args:", args);
  await deployUUPSProxiedContract(hre, deployer, "SequencerInbox", args);
};

export default func;
func.tags = ["SequencerInbox", "L1", "Stage0"];
