import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { Web3Provider, ExternalProvider } from "@ethersproject/providers";
import { getProxyName } from "../utils";

// Configure L1Portal with deployed L2 contract addresses
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, network, companionNetworks } =
    hre;

  // Get L1 deployer signer
  const { deployer } = await getNamedAccounts();
  const provider = new Web3Provider(
    companionNetworks.l1.provider as unknown as ExternalProvider
  );
  const deployerSigner = await provider.getSigner(deployer);

  // Get L1Portal contract on L1
  // const l1PortalProxyAddress = (
  //   await companionNetworks.l1.deployments.get(getProxyName("L1Portal"))
  // ).address;

  const l1PortalProxyAddress = "0x13D69Cf7d6CE4218F646B759Dcf334D82c023d8e";

  const l1Portal = await ethers.getContractAt(
    "L1Portal",
    l1PortalProxyAddress,
    deployerSigner
  );
  console.log(l1Portal.signer.getAddress());

  // Set L2Portal address on L1Portal
  const l2PortalProxyAddress = (await deployments.get(getProxyName("L2Portal")))
    .address;
  const tx = await l1Portal.setL2PortalAddress(l2PortalProxyAddress);
  const result = await tx.wait();
  console.log("L1Portal.setL2PortalAddress result:", result);
};

export default func;
func.tags = ["L2Portal", "L2", "postdeploy"];
