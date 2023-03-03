import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { Manifest } from "@openzeppelin/upgrades-core";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, upgrades } = hre;
  const { save } = deployments;
  const { deployer, sequencer } = await getNamedAccounts();
  console.log(sequencer);
  const deployerSigner = await ethers.getSigner(deployer);

  const L1Oracle = await ethers.getContractFactory("L1Oracle", deployerSigner);
  const l1Oracle = await upgrades.deployProxy(L1Oracle, [sequencer], {
    initializer: "initialize",
    timeout: 0,
    kind: "uups",
  });

  await l1Oracle.deployed();
  console.log("L1Oracle Proxy:", l1Oracle.address);
  console.log(
    "L1Oracle Implementation Address",
    await upgrades.erc1967.getImplementationAddress(l1Oracle.address)
  );
  console.log(
    "L1Oracle Admin Address",
    await upgrades.erc1967.getAdminAddress(l1Oracle.address)
  );

  const artifact = await deployments.getExtendedArtifact("L1Oracle");
  const proxyDeployments = {
    address: l1Oracle.address,
    ...artifact,
  };
  await save("L1Oracle", proxyDeployments);
};

export default func;
func.tags = ["L1Oracle", "L2"];
