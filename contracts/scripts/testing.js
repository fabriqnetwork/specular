/* eslint-disable no-undef */

// eslint-disable-next-line @typescript-eslint/no-var-requires
const rollupAbi = require("../abi/src/Rollup.sol/Rollup.json");
// eslint-disable-next-line @typescript-eslint/no-var-requires
const sequencerInboxAbi = require("../abi/src/SequencerInbox.sol/SequencerInbox.json");
// eslint-disable-next-line @typescript-eslint/no-var-requires
//const { ethers } = require("ethers");

// eslint-disable-next-line @typescript-eslint/no-var-requires
const { ethers } = require("hardhat");
// eslint-disable-next-line @typescript-eslint/no-var-requires
const { Wallet, utils } = require("ethers");

// eslint-disable-next-line @typescript-eslint/no-var-requires
const { spawn } = require("child_process");
// eslint-disable-next-line @typescript-eslint/no-var-requires
const { exec } = require("child_process");
// eslint-disable-next-line @typescript-eslint/no-var-requires
const net = require("net");
// eslint-disable-next-line @typescript-eslint/no-var-requires
const fs = require("fs");

const sequencerScriptPath =
  "../../clients/geth/specular/sbin/start_sequencer.sh";

const sequencerPrivateKeyPath =
  "../../clients/geth/specular/data/keys/sequencer.prv";

async function startL1() {
  const l1Process = spawn("npx", ["hardhat", "node"]);

  await new Promise((resolve, reject) => {
    l1Process.stdout.on("data", (data) => {
      console.log("L1 data: ", data.toString());
      if (data.toString().includes("Listening")) {
        console.log("L1 is ready...");
        resolve();
      }
    });

    l1Process.stderr.on("data", (data) => {
      reject(`Error starting L1 node: ${data.toString()}`);
    });
  });
}

async function checkL2Status() {
  return new Promise((resolve, reject) => {
    const client = new net.Socket();
    client.connect(4011, "127.0.0.1", () => {
      console.log("L2 is running on port 4011");
      resolve();
      client.end();
    });
    client.on("error", (err) => {
      console.log(`Error checking L2 status: ${err}`);
      reject(err);
    });
  });
}

async function startL2() {
  const options = {
    hostname: "localhost",
    port: 4011,
    path: "/",
    method: "HEAD",
  };

  const child = exec(`bash ${sequencerScriptPath}`);

  // Wait for the L2 node to print "READY" to stdout
  await new Promise((resolve, reject) => {
    child.stdout.on("data", (data) => {
      if (data.toString().includes("READY")) {
        console.log("L2 is ready...");
        resolve();
      }
    });

    child.stderr.on("data", (data) => {
      reject(`Error starting L2 node: ${data.toString()}`);
    });
  });

  await checkL2Status();
  console.log("L2 is up and running");
}

// Deploy Contracts on L1
async function deployContractsOnL1() {
  const provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:8545"
  );

  const sequencerPrivateKey = fs.readFileSync(sequencerPrivateKeyPath, "utf8");
  console.log(sequencerPrivateKey);

  const wallet = new Wallet(sequencerPrivateKey, provider);

  const signer = wallet.connect(provider);

  const rollupFactory = await ethers.getContractFactory("Rollup");
  const sequencerInboxFactory = await ethers.getContractFactory(
    "SequencerInbox"
  );

  const sequencerInbox = await sequencerInboxFactory.deploy();
  const rollup = await rollupFactory.deploy();

  return [sequencerInbox.address, rollup.address];
}

// Test txs
async function testTxs(
  sequencerInboxAddress,
  rollupAddress,
  fromAddress,
  toAddress,
  value
) {
  const provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:4011"
  );

  const sequencerPrivateKey = fs.readFileSync(sequencerPrivateKeyPath, "utf8");

  const signer = new Wallet(sequencerPrivateKey, provider);

  let nonce = await provider.getTransactionCount(signer.address);

  const sequencerContractAddress = sequencerInboxAddress;
  const rollupContractAddress = rollupAddress;

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

  // Check AppendTx event
  const appendTxLogs = await sequencerContract.queryFilter(
    appendTxFilter,
    txReceipt.blockHash
  );

  // Check Assertion creation
  const assertionCreatedLogs = await rollupContract.queryFilter(
    assertionCreatedFilter,
    txReceipt.blockHash
  );

  // Check Assertion confirmation
  const assertionConfirmedLogs = await rollupContract.queryFilter(
    assertionConfirmedFilter,
    txReceipt.blockHash
  );
}

// Send multiple Txs
async function sendMultipleTransactions() {
  await startL1();
  await startL2();

  let [sequencerInboxAddress, rollupAddress] = await deployContractsOnL1();

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
