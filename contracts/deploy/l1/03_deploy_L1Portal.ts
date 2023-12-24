import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, upgrades } = hre;
  const { deployer } = await getNamedAccounts();
  const rollupProxyAddress = (await deployments.get(getProxyName("Rollup")))
    .address;
  await deployUUPSProxiedContract(hre, deployer, "L1Portal", [
    rollupProxyAddress,
  ]);
};

export default func;
func.tags = ["L1Portal", "L1", "Stage0"];
