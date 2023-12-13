import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getStorageKey,
  getWithdrawalProof,
  delay,
} from "../utils";
import { BigNumber } from "ethers";

async function main() {
  const {
    l1Bridger,
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

  const initEvent = l2Portal.interface.parseLog(txWithLogs.logs[1]);
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

  let assertionWasCreated = false;
  let assertionId: number | undefined = undefined;
  let lastConfirmedBlockNum: number | undefined = undefined;

  inbox.on(
    inbox.filters.TxBatchAppended(),
    (event) => {
      console.log(`TxBatchAppended blockNum: ${event.blockNumber}`)
    }
  );

  rollup.on(rollup.filters.AssertionConfirmed(), async (id: BigNumber) => {
    if (assertionWasCreated) {
      assertionId = id.toNumber();
      const assertion = await rollup.getAssertion(assertionId);
      lastConfirmedBlockNum = assertion.blockNum.toNumber();
      console.log("AssertionConfirmed", "id", assertionId, "blockNum", lastConfirmedBlockNum)
    }
  });

  rollup.on(rollup.filters.AssertionCreated(), () => {
    console.log("AssertionCreated")
    assertionWasCreated = true;
  });

  l1StandardBridge.on(
    l1StandardBridge.filters.ETHBridgeFinalized(),
    async (from, to, amount, data) => {
      console.log({ msg: "ETHBridgeFinalized", from, to, amount, data });
    }
  );

  console.log("Waiting for assertion to be confirmed...");
  while (!assertionId || !lastConfirmedBlockNum || lastConfirmedBlockNum < blockNumber) {
    await delay(500);
  }

  const { accountProof, storageProof } = await getWithdrawalProof(
    l2Portal.address,
    initEvent.args.withdrawalHash
  );

  const finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
    crossDomainMessage,
    assertionId,
    accountProof,
    storageProof
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
