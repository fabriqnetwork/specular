import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
// import { Manifest } from "@openzeppelin/upgrades-core";
import { getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, upgrades } = hre;
  const { save } = deployments;
  const { deployer } = await getNamedAccounts();
  const deployerSigner = await ethers.getSigner(deployer);

  const rollupProxyAddress = (await deployments.get(getProxyName("Rollup")))
    .address;

  const L1Portal = await ethers.getContractFactory("L1Portal", deployerSigner);
  const l1Portal = await upgrades.deployProxy(L1Portal, [rollupProxyAddress], {
    initializer: "initialize",
    timeout: 0,
    kind: "uups",
  });

  await l1Portal.deployed();
  console.log("L1Portal Proxy:", l1Portal.address);
  console.log(
    "L1Portal Implementation Address",
    await upgrades.erc1967.getImplementationAddress(l1Portal.address)
  );
  console.log(
    "L1Portal Admin Address",
    await upgrades.erc1967.getAdminAddress(l1Portal.address)
  );

  const artifact = await deployments.getExtendedArtifact("L1Portal");
  const proxyDeployments = {
    address: l1Portal.address,
    ...artifact,
  };

  await save(getProxyName("L1Portal"), proxyDeployments);
};

export default func;
func.tags = ["L1Portal", "L1", "Stage0"];
