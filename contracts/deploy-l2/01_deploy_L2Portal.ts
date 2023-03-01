import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { Manifest } from "@openzeppelin/upgrades-core";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const {
    deployments,
    getNamedAccounts,
    ethers,
    upgrades,
    network,
    companionNetworks,
  } = hre;
  const { save } = deployments;
  const { deployer } = await getNamedAccounts();
  const deployerSigner = await ethers.getSigner(deployer);
  const { provider } = network;

  const l1OracleProxyAddress = (await deployments.get("L1Oracle")).address;
  const l1PortalProxyAddress = (
    await companionNetworks.l1.deployments.get("L1Portal")
  ).address;

  const l2PortalArgs = [l1OracleProxyAddress, l1PortalProxyAddress];

  const L2Portal = await ethers.getContractFactory("L2Portal", deployer);
  const l2Portal = await upgrades.deployProxy(L2Portal, l2PortalArgs, {
    initializer: "initialize",
    timeout: 0,
    kind: "uups",
  });

  await l2Portal.deployed();
  console.log("L2Portal Proxy:", l2Portal.address);
  console.log(
    "L2Portal Implementation Address",
    await upgrades.erc1967.getImplementationAddress(l2Portal.address)
  );
  console.log(
    "L2Portal Admin Address",
    await upgrades.erc1967.getAdminAddress(l2Portal.address)
  );

  const artifact = await deployments.getExtendedArtifact("L2Portal");
  const proxyDeployments = {
    address: l2Portal.address,
    ...artifact,
  };
  await save("L2Portal", proxyDeployments);
};

export default func;
func.tags = ["L2Portal", "L2"];
