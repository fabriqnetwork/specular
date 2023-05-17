import { exec } from "child_process";
import util from "node:util";
import path from "node:path";
import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract, getProxyName } from "../utils";

const CLIENT_SBIN_DIR = "../../../clients/geth/specular/sbin";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  // Calculate initial VM hash
  const execPromise = util.promisify(exec);
  let initialVMHash = "";
  try {
    const { stdout } = await execPromise(
      path.join(CLIENT_SBIN_DIR, "export_genesis.sh")
    );
    initialVMHash = (JSON.parse(stdout).root || "") as string;
    if (!initialVMHash) {
      throw Error(
        `could not export genesis hash, root field not found\n${stdout}`
      );
    }
    console.log("initial VM hash:", initialVMHash);
  } catch (err) {
    throw Error(`could not export genesis hash${err}`);
  }

  const { deployments, getNamedAccounts } = hre;
  const { sequencer, deployer } = await getNamedAccounts();
  const sequencerInboxProxyAddress = (
    await deployments.get(getProxyName("SequencerInbox"))
  ).address;
  const verifierProxyAddress = (await deployments.get(getProxyName("Verifier")))
    .address;

  const args = [
    sequencer, // address _vault
    sequencerInboxProxyAddress, // address _sequencerInbox
    verifierProxyAddress, // address _verifier
    5, // uint256 _confirmationPeriod
    0, // uint256 _challengePeriod
    0, // uint256 _minimumAssertionPeriod
    0, // uint256 _baseStakeAmount
    0, // uint256 _initialAssertionID
    0, // uint256 _initialInboxSize
    "initialVMHash", // bytes32 _initialVMhash
  ];

  await deployUUPSProxiedContract(hre, deployer, "Rollup", args);
};

export default func;
func.tags = ["Rollup", "L1", "Stage0"];
