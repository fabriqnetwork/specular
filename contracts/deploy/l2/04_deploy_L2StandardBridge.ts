import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deployer } = await getNamedAccounts();

  const l2PortalProxyAddress = (await deployments.get(getProxyName("L2Portal")))
    .address;

  // TODO: read this from companion chain deployments
  const l1StandardBridgeAddress = "0xE7C2a73131dd48D8AC46dCD7Ab80C8cbeE5b410A";

  const args = [l2PortalProxyAddress, l1StandardBridgeAddress];

  await deployUUPSProxiedContract(hre, deployer, "L2StandardBridge", args);
};

export default func;
func.tags = ["L2StandardBridge", "L2", "postdeploy"];
