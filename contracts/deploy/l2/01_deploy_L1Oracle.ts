import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployer } = await hre.getNamedAccounts();
  await deployUUPSProxiedContract(hre, deployer, "L1Oracle", []);
};

export default func;
func.tags = ["L1Oracle", "L2"];
