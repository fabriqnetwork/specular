import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, upgrades } = hre;
  const { save } = deployments;
  const { sequencer, deployer } = await getNamedAccounts();
  console.log(deployer);

  const Inbox = await ethers.getContractFactory("SequencerInbox", deployer);
  const inbox = await upgrades.deployProxy(Inbox, [sequencer], {
    initializer: "initialize",
    timeout: 0,
    kind: "uups",
  });

  await inbox.deployed();
  console.log("inbox Proxy:", inbox.address);
  console.log(
    "inbox Implementation Address",
    await upgrades.erc1967.getImplementationAddress(inbox.address)
  );
  console.log(
    "inbox Admin Address",
    await upgrades.erc1967.getAdminAddress(inbox.address)
  );

  const artifact = await deployments.getExtendedArtifact("SequencerInbox");
  const proxyDeployments = {
    address: inbox.address,
    ...artifact,
  };
  await save("SequencerInbox", proxyDeployments);
};

export default func;
func.tags = ["SequencerInbox", "L1", "Stage0"];
