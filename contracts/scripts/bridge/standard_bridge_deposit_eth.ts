import { ethers } from "hardhat";
import { getSignersAndContracts, getStorageKey } from "./utils";

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
    version: 1,
    nonce: initEvent.args.nonce,
    sender: initEvent.args.sender,
    target: initEvent.args.target,
    value: initEvent.args.value,
    gasLimit: initEvent.args.gasLimit,
    data: initEvent.args.data,
  };

  const blockNumber = await l1Provider.getBlockNumber();
  const rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);
  await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);

  const proof = await l1Provider.send("eth_getProof", [
    l1Portal.address,
    [getStorageKey(initEvent.args.depositHash)],
    "latest",
  ]);
  const accountProof = proof.accountProof;
  const storageProof = proof.storageProof[0].proof;

  const finalizeTx = await l2Portal.finalizeDepositTransaction(
    crossDomainMessage,
    accountProof,
    storageProof
  );
  await finalizeTx.wait();

  const balanceEnd = await l2Bridger.getBalance();
  if (!balanceEnd.sub(balanceStart).eq(bridgeValue)) {
    throw "unexpected end balance";
  }

  console.log("bridging ETH was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
