import path from "node:path";
import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployer } = await getNamedAccounts();

  console.log("Deploying Rollup")
  await deployUUPSProxiedContract(hre, deployer, "Rollup", []);
};

export default func;
func.tags = ["Rollup", "L1", "Stage0"];
