import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployer } = await hre.getNamedAccounts();
  await deployUUPSProxiedContract(hre, deployer, "Verifier", []);
};

export default func;
func.tags = ["Verifier", "L1", "Stage0"];
