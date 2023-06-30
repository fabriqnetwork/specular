import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, companionNetworks } = hre;
  const { deployer } = await getNamedAccounts();

  const l1OracleProxyAddress = (await deployments.get(getProxyName("L1Oracle")))
    .address;

  // const l1PortalProxyAddress = (
  //   await companionNetworks.l1.deployments.get(getProxyName("L1Portal"))
  // ).address;

  const l1PortalProxyAddress = "0x13D69Cf7d6CE4218F646B759Dcf334D82c023d8e";

  const args = [l1OracleProxyAddress, l1PortalProxyAddress];
  await deployUUPSProxiedContract(hre, deployer, "L2Portal", args);

  // fund l2 portal
  const l1PortalFactory = await ethers.getContractFactory("L2Portal");
  const l1Portal = await l1PortalFactory.attach(
    "0xBC9129Dc0487fc2E169941C75aABC539f208fb01"
  );
  const tx = await l1Portal.donateETH({ value: ethers.utils.parseUnits("10") });
  const r = await tx.wait();
  console.log({ msg: "funded l2 portal", receipt: r });
};

export default func;
func.tags = ["L2Portal", "L2"];
