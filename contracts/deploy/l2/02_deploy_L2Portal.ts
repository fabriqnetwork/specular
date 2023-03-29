import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract, getProxyName } from "../deploy-utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, companionNetworks } = hre;
  const { deployer } = await getNamedAccounts();
  const l1OracleProxyAddress = (await deployments.get(getProxyName("L1Oracle")))
    .address;
  const l1PortalProxyAddress = (
    await companionNetworks.l1.deployments.get(getProxyName("L1Portal"))
  ).address;
  const args = [l1OracleProxyAddress, l1PortalProxyAddress];
  await deployUUPSProxiedContract(hre, deployer, "L2Portal", args);
};

export default func;
func.tags = ["L2Portal", "L2"];
