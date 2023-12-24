import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deployer } = await getNamedAccounts();
  const rollupProxyAddress = (await deployments.get(getProxyName("Rollup")))
    .address;
  const args = [rollupProxyAddress];
  console.log("Deploying L1Portal with args:", args);
  await deployUUPSProxiedContract(hre, deployer, "L1Portal", args);
};

export default func;
func.tags = ["L1Portal", "L1", "Stage0"];
