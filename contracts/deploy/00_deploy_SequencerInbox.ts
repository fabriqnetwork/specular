import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, upgrades } = hre;
  const { deploy } = deployments;
  const { sequencer } = await getNamedAccounts();

  const Inbox = await ethers.getContractFactory("SequencerInbox");
  const inbox = await upgrades.deployProxy(Inbox, [sequencer], {
    initializer: "initialize",
    from: sequencer,
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
};

export default func;
func.tags = ["SequencerInbox"];
