import rollupJson from "../deployments/localhost/Rollup.json";
import sequencerInboxJson from "../deployments/localhost/SequencerInbox.json";
import { Wallet, utils, ethers, BigNumber } from "ethers";
import assert from "assert";
import fs from "fs";

const ROOT_DIR = __dirname + "/../../";

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

// Test tx flow
async function testTxs(toAddress: string, value: BigNumber) {
  // Signer is sequencer
  const sequencerPrivateKey = fs.readFileSync(sequencerPrivateKeyPath, "utf8");
  const signer = new Wallet(sequencerPrivateKey, l2Provider);
  const nonce = await l2Provider.getTransactionCount(signer.address);

  const sequencerContractAddress = sequencerInboxJson.address;
  const rollupContractAddress = rollupJson.address;

  // Contracts
  const sequencerContract = new ethers.Contract(
    sequencerContractAddress,
    sequencerInboxJson.abi,
    l1Provider
  );
  const rollupContract = new ethers.Contract(
    rollupContractAddress,
    rollupJson.abi,
    l1Provider
  );

  // Event filters
  const appendTxFilter = sequencerContract.filters.TxBatchAppended();
  const assertionCreatedFilter = rollupContract.filters.AssertionCreated();
  const assertionConfirmedFilter = rollupContract.filters.AssertionConfirmed();

  // Tx
  const txData = {
    to: toAddress,
    value: value,
    nonce: nonce,
  };

  // Send Tx to L2
  const txResponse = await signer.sendTransaction(txData);
  await txResponse.wait();

  // Check Tx added to L2
  const txReceipt = await l2Provider.getTransactionReceipt(txResponse.hash);
  assert(txReceipt, "No tx on L2 blockchain");

  // Check AppendTx event
  const appendTxLogs = await sequencerContract.queryFilter(appendTxFilter);
  assert(appendTxLogs.length > 0, "No appended txs");

  // Check Assertion creation
  const assertionCreatedLogs = await rollupContract.queryFilter(
    assertionCreatedFilter
  );
  assert(assertionCreatedLogs.length > 0, "No created assertions");

  // Check Assertion confirmation
  const assertionConfirmedLogs = await rollupContract.queryFilter(
    assertionConfirmedFilter
  );
  assert(assertionConfirmedLogs.length > 0, "No confirmed assertions");
}

// Send multiple Txs
async function sendMultipleTxs() {
  const validatorPrivateKey = fs.readFileSync(validatorPrivateKeyPath, "utf8");
  const validatorSigner = new Wallet(validatorPrivateKey, l2Provider);

  for (let i = 0; i < 1; i++) {
    const res = await testTxs(validatorSigner.address, utils.parseEther("0.1"));
    console.log("Done sending i = ", i);
  }
}

sendMultipleTxs();
