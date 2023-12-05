import { ethers } from "hardhat";
import { BigNumber } from "ethers";
import { formatEther } from "ethers/lib/utils";
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
    // This should not be required thanks to magi
    //l1Oracle,
  } = await getSignersAndContracts();

  const balanceStart: BigNumber = await l2Bridger.getBalance();
  const bridgeValue: BigNumber = ethers.utils.parseEther("0.1");

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

  console.log("Initial block", { blockNumber, stateRoot });

  const { accountProof, storageProof } = await getDepositProof(
    l1Portal.address,
    initEvent.args.depositHash
  );
  // This should not be required thanks to magi
  //await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);

  blockNumber = await l1Provider.getBlockNumber();
  rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);

  // This should not be required thanks to magi
  //await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);
  l1Provider.send()

  // TODO: make this a proper event driven
  await delay(8000);

  console.log("After deposit block", { blockNumber, stateRoot });

  const finalizeTx = await l2Portal.finalizeDepositTransaction(
    crossDomainMessage,
    accountProof,
    storageProof
  );
  await finalizeTx.wait();

  // balanceStart <- balance on L2 before bridging
  // balanceEnd <- balance on L2 after bridging
  // bridgeValue <- the value we're transfering
  // Expected: balanceEnd - balanceStart - bridgeValue ~== 0
  const balanceEnd: BigNumber = await l2Bridger.getBalance();
  const actualDiff: BigNumber = balanceEnd.sub(balanceStart).sub(bridgeValue).abs();
  const acceptableMargin: BigNumber = ethers.utils.parseEther("0.0001");

  if (!actualDiff.lt(acceptableMargin)) {
    const situation = {
      balanceStart: formatEther(balanceStart),
      bridgeValue: formatEther(bridgeValue),
      balanceEnd: formatEther(balanceEnd),
      actualDiff: formatEther(actualDiff),
      acceptableMargin: formatEther(acceptableMargin),
    }
    console.log(situation);
    throw "value after bridging is not as expected, actualDiff is expected to be close to zero";
  }

  console.log("bridging ETH was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
