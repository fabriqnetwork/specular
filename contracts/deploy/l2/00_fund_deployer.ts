import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

// Fund L2 deployer with ETH from pre-allocated L2 sequencer account
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, ethers } = hre;
  const { sequencer, deployer } = await getNamedAccounts();

  // If deployer already has ETH, skip
  const balance = await ethers.provider.getBalance(deployer);
  if (balance.toNumber() > 0) {
    return;
  }

  // Transfer ETH from sequencer to deployer
  const sequencerSigner = await ethers.provider.getSigner(sequencer);
  const amount = ethers.utils.parseEther("10.0");
  const transferTx = {
    to: deployer,
    value: amount,
  };
  const tx = await sequencerSigner.sendTransaction(transferTx);
  await tx.wait();
};

export default func;
func.tags = ["L2", "predeploy"];
