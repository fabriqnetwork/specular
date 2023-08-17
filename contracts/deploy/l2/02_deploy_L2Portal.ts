import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract, getProxyName } from "../utils";

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

  // fund l2 portal
  const l2PortalFactory = await ethers.getContractFactory("L2Portal");
  const l2Portal = await l2PortalFactory.attach(
    (
      await deployments.get(getProxyName("L2Portal"))
    ).address
  );

  const value = ethers.utils.parseUnits(process.env.L2_PORTAL_FUNDING_ETH);
  const tx = await l2Portal.donateETH({ value });
  const r = await tx.wait();
  console.log({ msg: "funded l2 portal", receipt: r });
};

export default func;
func.tags = ["L2Portal", "L2"];
