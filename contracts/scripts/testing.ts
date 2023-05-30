import rollupJson from "../abi/src/IRollup.sol/IRollup.json";
import sequencerInboxJson from "../abi/src/SequencerInbox.sol/SequencerInbox.json";
import { Wallet, utils, ethers, BigNumber } from "ethers";
import assert from "assert";
import fs from "fs";

const ROOT_DIR = __dirname + "/../../";

function delay(seconds: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, seconds * 1000));
}

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

// Send a tx
async function sendTx(sequencerSigner: any, toAddress: any, value: number) {
  const nonce = await l2Provider.getTransactionCount(sequencerSigner.address);

  const txData = {
    to: toAddress,
    value: value,
    nonce: nonce,
  };

  const txResponse = await sequencerSigner.sendTransaction(txData);
  await txResponse.wait();

  const txReceipt = await l2Provider.getTransactionReceipt(txResponse.hash);
  assert(txReceipt, "No tx on L2 blockchain");

  return txResponse;
}

// Check logs
async function checkLogs(name: string, contract: any, filter: any) {
  const logs = await contract.queryFilter(filter);
  assert(logs.length > 0, `No matching logs found for ${name}`);
}

// test Tx
async function testTx(
  sequencerSigner: any,
  validatorSigner: any,
  sequencerContract: any,
  rollupContract: any,
  appendTxFilter: any,
  assertionCreatedFilter: any,
  assertionConfirmedFilter: any,
  toAddress: any,
  value: any
) {
  const txResponse = await sendTx(sequencerSigner, toAddress, value);

  await delay(60);
  await checkLogs("appendTxFilter", sequencerContract, appendTxFilter);
  await checkLogs(
    "assertionCreatedFilter",
    rollupContract,
    assertionCreatedFilter
  );
  await checkLogs(
    "assertionConfirmedFilter",
    rollupContract,
    assertionConfirmedFilter
  );

  return txResponse;
}

// New Test tx flow
async function testTxs(toAddress: string, value: BigNumber) {
  const { sequencerSigner, validatorSigner } = await setupSigners(
    sequencerPrivateKeyPath,
    validatorPrivateKeyPath
  );

  // Contract Address Hardcoded
  const sequencerContractAddress = "0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1";
  const rollupContractAddress = "0xF6168876932289D073567f347121A267095f3DD6";

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

  for (let i = 0; i < 1; i++) {
    const res = await testTx(
      sequencerSigner,
      validatorSigner,
      sequencerContract,
      rollupContract,
      appendTxFilter,
      assertionCreatedFilter,
      assertionConfirmedFilter,
      toAddress,
      value
    );
  }
}

// Send Txs
async function sendTxs() {
  const validatorPrivateKey = fs.readFileSync(validatorPrivateKeyPath, "utf8");
  const validatorSigner = new Wallet(validatorPrivateKey, l2Provider);
  await testTxs(validatorSigner.address, utils.parseEther("0.1"));
}

sendTxs()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
