import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { sequencer, deployer } = await hre.getNamedAccounts();
  const bridger = "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199"
  await deployUUPSProxiedContract(hre, deployer, "L1Oracle", [bridger]);
};

export default func;
func.tags = ["L1Oracle", "L2"];
