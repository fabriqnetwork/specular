import path from "node:path";
import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract, getProxyName } from "../utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts } = hre;
  const { deployer } = await getNamedAccounts();
  const sequencerInboxProxyAddress = (
    await deployments.get(getProxyName("SequencerInbox"))
  ).address;
  const verifierProxyAddress = (await deployments.get(getProxyName("Verifier")))
    .address;

  const config = {
    vault: process.env.DEPLOYER_ADDRESS,
    daProvider: sequencerInboxProxyAddress,
    verifier: verifierProxyAddress,
    confirmationPeriod: 12, // TODO(#302): move to config
    challengePeriod: 0,
    minimumAssertionPeriod: 0,
    baseStakeAmount: 0,
    validators: [process.env.VALIDATOR_ADDRESS]
  };

  const args = [config];

  console.log("Deploying Rollup with args:", { args })
  await deployUUPSProxiedContract(hre, deployer, "Rollup", args);
};

export default func;
func.tags = ["Rollup", "L1", "Stage0"];
