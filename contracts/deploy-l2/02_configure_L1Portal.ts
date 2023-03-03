import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { Manifest } from "@openzeppelin/upgrades-core";
import { Web3Provider, ExternalProvider } from "@ethersproject/providers";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, network, companionNetworks } =
    hre;
  const { deployer } = await getNamedAccounts();
  const provider = new Web3Provider(
    companionNetworks.l1.provider as unknown as ExternalProvider
  );
  const deployerSigner = await provider.getSigner(deployer);

  const l1PortalProxyAddress = (
    await companionNetworks.l1.deployments.get("L1Portal")
  ).address;
  const l1Portal = await ethers.getContractAt(
    "L1Portal",
    l1PortalProxyAddress,
    deployerSigner
  );
  console.log(l1Portal.signer.getAddress());
  const l2PortalProxyAddress = (await deployments.get("L2Portal")).address;
  const tx = await l1Portal.setL2PortalAddress(l2PortalProxyAddress);
  const result = await tx.wait();
  console.log("L1Portal.setL2PortalAddress result:", result);
};

export default func;
func.tags = ["L2Portal", "L2", "postdeploy"];
