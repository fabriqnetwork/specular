import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deployer } = await getNamedAccounts();
  const l1PortalProxyAddress = (await deployments.get(getProxyName("L1Portal")))
    .address;
  const args = [l1PortalProxyAddress];
  console.log("Deploying L1StandardBridge with args:", args);
  await deployUUPSProxiedContract(hre, deployer, "L1StandardBridge", args);
};

export default func;
func.tags = ["L1StandardBridge", "L1", "Stage0"];
