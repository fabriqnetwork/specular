import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getStorageKey,
  getDepositProof,
  getWithdrawalProof,
  delay,
  getLastBlockNumber,
  deployTokenPair,
} from "../utils";

async function main() {
  const {
    l1Provider,
    l2Provider,
    l1Bridger,
    l2Bridger,
    l2Relayer,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2StandardBridge,
    l1Oracle,
    rollup,
    inbox,
  } = await getSignersAndContracts();

  const { l1Token, l2Token } = await deployTokenPair(l1Bridger, l2Relayer);
  console.log("\tdeployed tokens...");

  const l1BalanceStart = await l1Token.balanceOf(l1Bridger.address);

  const approveTx = await l1Token.approve(
    l1StandardBridge.address,
    l1BalanceStart
  );
  await approveTx.wait();

  const depositTx = await l1StandardBridge.bridgeERC20(
    l1Token.address,
    l2Token.address,
    l1BalanceStart,
    200_000,
    []
  );
  const depositTxWithLogs = await depositTx.wait();
  const l1BalanceEnd = await l1Token.balanceOf(l1Bridger.address);

  const depositEvent = await l1Portal.interface.parseLog(
    depositTxWithLogs.logs[3]
  );
  const depositMessage = {
    version: 0,
    nonce: depositEvent.args.nonce,
    sender: depositEvent.args.sender,
    target: depositEvent.args.target,
    value: depositEvent.args.value,
    gasLimit: depositEvent.args.gasLimit,
    data: depositEvent.args.data,
  };

  const depositBlockNumber = await l1Provider.getBlockNumber();
  const rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(depositBlockNumber),
    false, // We only want the block header
  ]);
  const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);

  const depositProof = await getDepositProof(
    l1Portal.address,
    depositEvent.args.depositHash
  );
  await l1Oracle.setL1OracleValues(depositBlockNumber, stateRoot, 0);

  const tx = await l2Portal.finalizeDepositTransaction(
    depositMessage,
    depositProof.accountProof,
    depositProof.storageProof
  );

  await tx.wait();

  const l2BalanceEnd = await l2Token.balanceOf(l2Bridger.address);

  if (!l1BalanceEnd.eq(0) || !l2BalanceEnd.eq(l1BalanceStart)) {
    throw "unexpected end balance";
  }
  console.log("\tdeposited token...");

  const withdrawalTx = await l2StandardBridge.bridgeERC20(
    l2Token.address,
    l1Token.address,
    l2BalanceEnd,
    200_000,
    []
  );
  const txWithLogs = await withdrawalTx.wait();

  const l2BalanceEmpty = await l2Token.balanceOf(l2Bridger.address);
  if (!l2BalanceEmpty.eq(0)) {
    throw "unexpected L2 balance";
  }

  const withdrawalEvent = await l2Portal.interface.parseLog(txWithLogs.logs[3]);
  const crossDomainMessage = {
    version: 0,
    nonce: withdrawalEvent.args.nonce,
    sender: withdrawalEvent.args.sender,
    target: withdrawalEvent.args.target,
    value: withdrawalEvent.args.value,
    gasLimit: withdrawalEvent.args.gasLimit,
    data: withdrawalEvent.args.data,
  };

  const blockNumber = txWithLogs.blockNumber;

  let lastConfirmedBlockNumber = 0;
  let assertionWasCreated = false;
  let assertionId;

  inbox.on(
    inbox.filters.TxBatchAppended(),
    async (batchNumber, previousInboxSize, inboxSize, event) => {
      const tx = await event.getTransaction();
      lastConfirmedBlockNumber = await getLastBlockNumber(tx.data);
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

  console.log("\twaiting for assertion to be confirmed...");

  while (lastConfirmedBlockNumber < blockNumber || !assertionId) {
    await delay(500);
  }

  const withdrawalProof = await getWithdrawalProof(
    l2Portal.address,
    withdrawalEvent.args.withdrawalHash,
    `0x${lastConfirmedBlockNumber.toString(16)}`
  );

  const finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
    crossDomainMessage,
    assertionId,
    withdrawalProof.accountProof,
    withdrawalProof.storageProof
  );
  await finalizeTx.wait();

  const finalBalance = await l1Token.balanceOf(l1Bridger.address);
  if (!finalBalance.eq(l1BalanceStart)) {
    throw "unexpected end balance";
  }

  console.log("withdrawing ERC20 was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
