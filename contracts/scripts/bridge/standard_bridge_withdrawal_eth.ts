import { ethers } from "hardhat";
import { getSignersAndContracts, getStorageKey } from "./utils";

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function getLastBlockNumber(data) {
  const iface = new ethers.utils.Interface([
    "function appendTxBatch(uint256[],uint256[],uint256,bytes)",
  ]);
  const decoded = iface.decodeFunctionData(
    "appendTxBatch(uint256[],uint256[],uint256,bytes)",
    data
  );
  const contexts: BigNumber[] = decoded[0];
  const firstL2BlockNumber = decoded[2];
  const lastL2BlockNumber =
    contexts.length / 2 + firstL2BlockNumber.toNumber() - 1;
  return lastL2BlockNumber;
}

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
  let assertionWasCreated = false;
  let assertionId;

  inbox.on(
    inbox.filters.TxBatchAppended(),
    async (batchNumber, previousInboxSize, inboxSize, event) => {
      const tx = await event.getTransaction();
      lastConfirmedBlockNumber = getLastBlockNumber(tx.data, inbox);
    }
  );

  rollup.on(rollup.filters.AssertionConfirmed(), (id) => {
    if (assertionWasCreated) {
      assertionId = id;
    }
  });

  rollup.on(rollup.filters.AssertionCreated(), () => {
    assertionWasCreated = true;
  });

  l1StandardBridge.on(
    l1StandardBridge.filters.ETHBridgeFinalized(),
    async (from, to, amount, data) => {
      console.log({ msg: "ETHBridgeFinalized", from, to, amount, data });
    }
  );

  console.log("\twaiting for assertion to be confirmed...");
  while (lastConfirmedBlockNumber < blockNumber || !assertionId) {
    await delay(500);
  }

  const proof = await l2Provider.send("eth_getProof", [
    l2Portal.address,
    [getStorageKey(initEvent.args.withdrawalHash)],
    ethers.utils.hexValue(lastConfirmedBlockNumber),
  ]);
  const accountProof = proof.accountProof;
  const storageProof = proof.storageProof[0].proof;

  const finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
    crossDomainMessage,
    assertionId,
    accountProof,
    storageProof
  );
  await finalizeTx.wait();

  const balanceEnd = await l1Bridger.getBalance();
  console.log({
    start: ethers.utils.formatUnits(balanceStart),
    end: ethers.utils.formatUnits(balanceEnd),
    diff: ethers.utils.formatUnits(balanceEnd.sub(balanceStart)),
  });
  if (!balanceEnd.sub(balanceStart).eq(bridgeValue)) {
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
