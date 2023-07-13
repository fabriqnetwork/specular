import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deployer } = await getNamedAccounts();

  const l1PortalProxyAddress = (await deployments.get(getProxyName("L1Portal")))
    .address;

  // TODO: this should be a pre-deploy address
  const l2StandardBridgeAddress = "0xF6168876932289D073567f347121A267095f3DD6";
  const args = [l1PortalProxyAddress, l2StandardBridgeAddress];

  await deployUUPSProxiedContract(hre, deployer, "L1StandardBridge", args);
};

export default func;
func.tags = ["L1StandardBridge", "L1", "Stage0"];
