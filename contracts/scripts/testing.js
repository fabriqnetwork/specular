// Import ABIs
// eslint-disable-next-line @typescript-eslint/no-var-requires
const rollupAbi = require("../abi/src/Rollup.sol/Rollup.json");
// eslint-disable-next-line @typescript-eslint/no-var-requires
const sequencerInboxAbi = require("../abi/src/SequencerInbox.sol/SequencerInbox.json");
// eslint-disable-next-line @typescript-eslint/no-var-requires
const { Wallet, utils, ethers } = require("ethers");

async function testTxs(fromAddress, toAddress, value) {
  const provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:4011"
  );

  const privateKey =
    "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80";

  const signer = new Wallet(privateKey, provider);

  let nonce = await provider.getTransactionCount(signer.address);

  const sequencerContractAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";
  const rollupContractAddress = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";

  // Contracts
  const sequencerContract = new ethers.Contract(
    sequencerContractAddress,
    sequencerInboxAbi,
    provider
  );
  const rollupContract = new ethers.Contract(
    rollupContractAddress,
    rollupAbi,
    provider
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
  const txReceipt = await provider.getTransactionReceipt(txResponse.hash);
  console.log("Transaction added to L2 blockchain: ", txReceipt);

  // Check AppendTx event
  const appendTxLogs = await sequencerContract.queryFilter(
    appendTxFilter,
    txReceipt.blockHash
  );

  console.log("AppendTx event emitted:", appendTxLogs);

  // Check Assertion creation
  const assertionCreatedLogs = await rollupContract.queryFilter(
    assertionCreatedFilter,
    txReceipt.blockHash
  );
  console.log("Assertion created:", assertionCreatedLogs);

  // Check Assertion confirmation
  const assertionConfirmedLogs = await rollupContract.queryFilter(
    assertionConfirmedFilter,
    txReceipt.blockHash
  );
  console.log("Assertion confirmed:", assertionConfirmedLogs);
}

// Send multiple Txs
async function sendMultipleTransactions() {
  for (let i = 0; i < 1; i++) {
    const res = await testTxs(
      "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
      "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      utils.parseEther("0.1")
    );
    console.log("Done sending i = ", i);
  }
}

sendMultipleTransactions();
console.log("Done with e2e testing script");
