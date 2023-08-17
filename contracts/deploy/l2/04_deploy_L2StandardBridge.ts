import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deployer } = await getNamedAccounts();

  const l2PortalProxyAddress = (await deployments.get(getProxyName("L2Portal")))
    .address;

  const l1StandardBridgeAddress = (
    await companionNetworks.l1.deployments.get(getProxyName("L1StandardBridge"))
  ).address;

  const args = [l2PortalProxyAddress, l1StandardBridgeAddress];

  await deployUUPSProxiedContract(hre, deployer, "L2StandardBridge", args);
};

export default func;
func.tags = ["L2StandardBridge", "L2", "postdeploy"];
