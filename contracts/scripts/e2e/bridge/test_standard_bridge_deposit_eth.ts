import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getStorageKey,
  getDepositProof,
  delay,
} from "../utils";

async function main() {
  const {
    l1Provider,
    l2Bridger,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l1Oracle,
  } = await getSignersAndContracts();

  const balanceStart = await l2Bridger.getBalance();
  const bridgeValue = ethers.utils.parseEther("0.1");

  const bridgeTx = await l1StandardBridge.bridgeETH(200_000, [], {
    value: bridgeValue,
  });
  const txWithLogs = await bridgeTx.wait();

  const initEvent = await l1Portal.interface.parseLog(txWithLogs.logs[1]);
  const crossDomainMessage = {
    version: 0,
    nonce: initEvent.args.nonce,
    sender: initEvent.args.sender,
    target: initEvent.args.target,
    value: initEvent.args.value,
    gasLimit: initEvent.args.gasLimit,
    data: initEvent.args.data,
  };

  let blockNumber = await l1Provider.getBlockNumber();
  let rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  let stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);
  console.log({ blockNumber, stateRoot });

  const { accountProof, storageProof } = await getDepositProof(
    l1Portal.address,
    initEvent.args.depositHash
  );
  await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);

  blockNumber = await l1Provider.getBlockNumber();
  rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);
  console.log({ blockNumber, stateRoot });
  await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);

  const finalizeTx = await l2Portal.finalizeDepositTransaction(
    crossDomainMessage,
    accountProof,
    storageProof
  );
  await finalizeTx.wait();

  const balanceEnd = await l2Bridger.getBalance();
  const balanceDiff = balanceEnd.sub(balanceStart);
  const error = ethers.utils.parseEther("0.001");

  if (bridgeValue.sub(balanceDiff).gt(error)) {
    console.log("bridging ETH was successful");
    throw "unexpected end balance";
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
