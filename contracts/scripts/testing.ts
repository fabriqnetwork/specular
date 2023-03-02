import rollupAbi from "../abi/src/Rollup.sol/Rollup.json";
import sequencerInboxAbi from "../abi/src/SequencerInbox.sol/SequencerInbox.json";
import { Wallet, utils, ethers, BigNumber } from "ethers";
import assert from "assert";
import fs from "fs";

const sequencerPrivateKeyPath =
  "../../clients/geth/specular/data/keys/sequencer.prv";

// Test tx flow
async function testTxs(toAddress: string, value: BigNumber) {
  const l2Provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:4011"
  );

  const l1Provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:8545"
  );

  // Signer is sequencer
  const sequencerPrivateKey = fs.readFileSync(sequencerPrivateKeyPath, "utf8");
  const signer = new Wallet(sequencerPrivateKey, l2Provider);
  const nonce = await l2Provider.getTransactionCount(signer.address);

  const sequencerContractAddress = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";
  const rollupContractAddress = "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707";

  // Contracts
  const sequencerContract = new ethers.Contract(
    sequencerContractAddress,
    sequencerInboxAbi,
    l1Provider
  );
  const rollupContract = new ethers.Contract(
    rollupContractAddress,
    rollupAbi,
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
  for (let i = 0; i < 1; i++) {
    const res = await testTxs(
      "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      utils.parseEther("0.1")
    );
    console.log("Done sending i = ", i);
  }
}

sendMultipleTxs();
