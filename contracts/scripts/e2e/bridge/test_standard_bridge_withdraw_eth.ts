import { Console } from "console";
import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getStorageKey,
  getWithdrawalProof,
  delay,
  getLastBlockNumber,
} from "../utils";

async function main() {
  const {
    l1Provider,
    l2Provider,
    l1Bridger,
    l2Relayer,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2StandardBridge,
    rollup,
    inbox,
  } = await getSignersAndContracts();

  const donateTx = await l1Portal.donateETH({
    value: ethers.utils.parseEther("1"),
  });
  await donateTx.wait();

  const balanceStart = await l1Bridger.getBalance();
  const bridgeValue = ethers.utils.parseEther("0.1");

  const bridgeTx = await l2StandardBridge.bridgeETH(200_000, [], {
    value: bridgeValue,
  });
  const txWithLogs = await bridgeTx.wait();

  const initEvent = await l2Portal.interface.parseLog(txWithLogs.logs[1]);
  const crossDomainMessage = {
    version: 0,
    nonce: initEvent.args.nonce,
    sender: initEvent.args.sender,
    target: initEvent.args.target,
    value: initEvent.args.value,
    gasLimit: initEvent.args.gasLimit,
    data: initEvent.args.data,
  };

  const blockNumber = txWithLogs.blockNumber;

  let lastConfirmedBlockNumber = 0;
  let assertionId;

  rollup.on(rollup.filters.AssertionConfirmed(), async (id: Number) => {
    assertionId = id;
    const assertion = await rollup.getAssertion(id)
    console.log({ assertion })
    lastConfirmedBlockNumber = assertion.blockNum.toNumber()
  });

  l1StandardBridge.on(
    l1StandardBridge.filters.ETHBridgeFinalized(),
    async (from, to, amount, data) => {
      console.log({ msg: "ETHBridgeFinalized", from, to, amount, data });
    }
  );

  console.log("\twaiting for assertion to be confirmed...");
  while (lastConfirmedBlockNumber < blockNumber || !assertionId) {
    console.log({ lastConfirmedBlockNumber, blockNumber, assertionId })
    await delay(500);
  }
  console.log({ lastConfirmedBlockNumber, blockNumber, assertionId })
  await delay(5000);

  const lastBlockNumber = (await l2Provider.getBlock("latest")).number

  console.log({ lastBlockNumber })

  for (let blockNumber = 0; blockNumber <= lastBlockNumber; blockNumber++) {
    let hexBlockNum = ethers.utils.hexlify(blockNumber)
    if (hexBlockNum.startsWith('0x0')) {
      hexBlockNum = '0x' + hexBlockNum.substr(3)
    }
    let rawBlock = await l2Provider.send("eth_getBlockByNumber", [
      hexBlockNum,
      false, // We only want the block header
    ]);
    const blockHash = rawBlock.hash;
    const stateRoot = rawBlock.stateRoot;

    const tmp = ethers.utils.concat([
      ethers.constants.HashZero,
      blockHash,
      stateRoot
    ])

    const stateCommitment = ethers.utils.keccak256(tmp)
    console.log({ blockNumber, blockHash, stateRoot, stateCommitment })
  }

  return

  let hexBlockNum = ethers.utils.hexlify(lastConfirmedBlockNumber + 1)
  if (hexBlockNum.startsWith('0x0')) {
    hexBlockNum = '0x' + hexBlockNum.substr(3)
  }
  const { accountProof, storageProof } = await getWithdrawalProof(
    l2Portal.address,
    initEvent.args.withdrawalHash,
    hexBlockNum
  );

  let rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    hexBlockNum,
    false, // We only want the block header
  ]);
  let stateRoot = l2Provider.formatter.hash(rawBlock.stateRoot);

  const tmp = ethers.utils.concat([
    ethers.constants.HashZero,
    stateRoot
  ])

  const rootWithVersion = ethers.utils.keccak256(tmp)
  console.log({ tmp, rootWithVersion })
  console.log({ stateRoot })
  const decode = ethers.utils.RLP.decode(storageProof[0])
  console.log({ storageProof })
  console.log({ decode })
  decode.push("0x00")
  const newProof = ethers.utils.RLP.encode(decode)
  const finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
    crossDomainMessage,
    assertionId,
    accountProof,
    newProof
  );
  await finalizeTx.wait();

  const balanceEnd = await l1Bridger.getBalance();
  const balanceDiff = balanceEnd.sub(balanceStart);
  const error = ethers.utils.parseEther("0.0001");

  if (bridgeValue.sub(balanceDiff).gt(error)) {
    throw "unexpected end balance";
  }

  console.log("withdrawing ETH was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
