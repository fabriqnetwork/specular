import { exec } from "child_process";
import util from "node:util";
import path from "node:path";
import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

import { deployUUPSProxiedContract, getProxyName } from "../utils";

const CONTRACTS_DIR = path.join(__dirname, "/../../")
require("dotenv").config({ path: path.join(CONTRACTS_DIR, ".genesis.env")});
const GENESIS_JSON = require(path.join(CONTRACTS_DIR, process.env.GENESIS_EXPORTED_HASH_PATH));

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  console.log({ GENESIS_JSON })
  const initialBlockHash = (GENESIS_JSON.hash || "") as string;
  if (!initialBlockHash) {
     throw Error(`blockHash not found\n$`);
  }
  console.log("initial blockHash:", initialBlockHash);

  const initialStateRoot = (GENESIS_JSON.stateRoot || "") as string;
  if (!initialStateRoot) {
     throw Error(`stateRoot not found\n$`);
  }
  console.log("initial stateRoot:", initialStateRoot);

  const { deployments, getNamedAccounts } = hre;
  const { sequencer, validator, deployer } = await getNamedAccounts();
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
    initialBlockHash, // bytes32 _initialBlockHash
    initialStateRoot, // bytes32 _initialStateRoot
    [sequencer, validator], // address[] calldata _validators
  ];

  await deployUUPSProxiedContract(hre, deployer, "Rollup", args);
};

export default func;
func.tags = ["Rollup", "L1", "Stage0"];
