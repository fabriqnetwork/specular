import { executeCommand } from "./docker_utils";
import path from "path";
import { Wallet, utils, ethers, BigNumber } from "ethers";
import assert from "assert";
import fs from "fs";
import rollupJson from "../deployments/localhost/Rollup.json";
import sequencerInboxJson from "../deployments/localhost/SequencerInbox.json";

const ROOT_DIR = path.join(__dirname, "../../");
const SPECULAR_DATADIR =
  process.env.SPECULAR_DATADIR || path.join(ROOT_DIR, "specular-datadir");

const sequencerPrivateKeyPath =
  ROOT_DIR + "clients/geth/specular/data/keys/sequencer.prv";
const validatorPrivateKeyPath =
  ROOT_DIR + "clients/geth/specular/data/keys/validator.prv";

const l2Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:4011"
);

const l1Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:8545"
);

// Setup signers
export async function setupSigners(
  sequencerPrivateKeyPath: string,
  validatorPrivateKeyPath: string
) {
  const sequencerPrivateKey = fs.readFileSync(sequencerPrivateKeyPath, "utf8");
  const sequencerSigner = new Wallet(sequencerPrivateKey, l2Provider);

  const validatorPrivateKey = fs.readFileSync(validatorPrivateKeyPath, "utf8");
  const validatorSigner = new Wallet(validatorPrivateKey, l2Provider);

  return {
    sequencerSigner,
    validatorSigner,
  };
}

// Initialize contracts and event filters
function initializeContracts(
  sequencerContractAddress: string,
  sequencerContractAbi: any,
  rollupContractAddress: string,
  rollupContractAbi: any
) {
  const sequencerContract = new ethers.Contract(
    sequencerContractAddress,
    sequencerContractAbi,
    l1Provider
  );
  const rollupContract = new ethers.Contract(
    rollupContractAddress,
    rollupContractAbi,
    l1Provider
  );

  const appendTxFilter = sequencerContract.filters.TxBatchAppended();
  const assertionCreatedFilter = rollupContract.filters.AssertionCreated();
  const assertionConfirmedFilter = rollupContract.filters.AssertionConfirmed();

  return {
    sequencerContract,
    rollupContract,
    appendTxFilter,
    assertionCreatedFilter,
    assertionConfirmedFilter,
  };
}

// Start L1 & L2
async function startChains() {
  console.log("In docker_start, this is ROOT_DIR: ", ROOT_DIR);
  console.log("In docker_start, this is SPECULAR_DATADIR: ", SPECULAR_DATADIR);
  const command = `docker-compose -f ${ROOT_DIR}docker/services/docker-compose-l1hardhat.yml up -d --build`;
  await executeCommand(command);
}

// Check logs
async function checkLogs(contract: any, filter: any) {
  const logs = await contract.queryFilter(filter);
  console.log("Logs: ", logs);
  assert(logs.length > 0, "No matching logs found");

  // Check for vmHash field in the logs
  const log = logs[0];
  assert(log.args && log.args.vmHash, "No vmHash field found in logs");
  console.log("vmHash: ", log.args.vmHash);
}

async function main() {
  // Start L1 & L2
  console.log("MAIN -> starting L1..");
  await startChains();

  console.log("MAIN -> done starting L1..");

  // Get the L1 contract
  const l1provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:8545"
  );

  const { sequencerSigner, validatorSigner } = await setupSigners(
    sequencerPrivateKeyPath,
    validatorPrivateKeyPath
  );

  const sequencerContractAddress = sequencerInboxJson.address;
  const rollupContractAddress = rollupJson.address;

  console.log("MAIN... sequencerContractAddress: ", sequencerContractAddress);

  const {
    sequencerContract,
    rollupContract,
    appendTxFilter,
    assertionCreatedFilter,
    assertionConfirmedFilter,
  } = initializeContracts(
    sequencerContractAddress,
    sequencerInboxJson.abi,
    rollupContractAddress,
    rollupJson.abi
  );

  console.log("MAIN... now checking logs: ");
  //   const l1wallet = new ethers.Wallet(
  //     process.env.L1_DEPLOYER_PRIVATE_KEY,
  //     l1provider
  //   );
  //   const L1Chain = require("../build/L1Chain.json");
  //   const l1Contract = new ethers.Contract(
  //     L1Chain.networks[31337].address,
  //     L1Chain.abi,
  //     l1wallet
  //   );

  // Wait for AssertionCreated event
  await checkLogs(rollupContract, assertionCreatedFilter);
  console.log("MAIN... done checking logs: ");
}

main();
