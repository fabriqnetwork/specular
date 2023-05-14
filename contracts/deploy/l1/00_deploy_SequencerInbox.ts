import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { sequencer, deployer } = await hre.getNamedAccounts();
  await deployUUPSProxiedContract(hre, deployer, "SequencerInbox", [sequencer]);
};

export default func;
func.tags = ["SequencerInbox", "L1", "Stage0"];
