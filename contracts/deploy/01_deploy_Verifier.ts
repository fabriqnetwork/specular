import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, upgrades } = hre;
  const { deploy, save } = deployments;
  const { sequencer, deployer } = await getNamedAccounts();
  const deployerSigner = await ethers.getSigner(deployer);

  const Verifier = await ethers.getContractFactory("Verifier", deployer);
  const verifier = await upgrades.deployProxy(Verifier, [], {
    initializer: "initialize",
    timeout: 0,
    kind: "uups",
  });

  await verifier.deployed();
  console.log("Verifier Proxy:", verifier.address);
  console.log(
    "Verifier Implementation Address",
    await upgrades.erc1967.getImplementationAddress(verifier.address)
  );
  console.log(
    "Verifier Admin Address",
    await upgrades.erc1967.getAdminAddress(verifier.address)
  );

  const artifact = await deployments.getExtendedArtifact("Verifier");
  const proxyDeployments = {
    address: verifier.address,
    ...artifact,
  };
  await save("Verifier", proxyDeployments);
};

export default func;
func.tags = ["Verifier"];
