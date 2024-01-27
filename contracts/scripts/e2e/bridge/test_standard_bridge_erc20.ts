import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getDepositProof,
  getWithdrawalProof,
  deployTokenPair,
  hexlifyBlockNum,
  waitUntilOracleBlock,
  waitUntilBlockConfirmed,
} from "../utils";

async function main() {
  const {
    l1Provider,
    l1Bridger,
    l2Bridger,
    l2Relayer,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2Provider,
    l2StandardBridge,
    l1Oracle,
    rollup,
  } = await getSignersAndContracts();

  const { l1Token, l2Token } = await deployTokenPair(l1Bridger, l2Relayer);
  console.log(`Deployed tokens ${l1Token.address}, ${l2Token.address}`);

  // TODO(#306): portal should be funded as part of pre-deploy pipeline
  await Promise.all([
    l1Portal.donateETH({ value: ethers.utils.parseEther("1") }),
    l2Portal.donateETH({ value: ethers.utils.parseEther("1") }),
  ]);

  const l1BalanceStart = await l1Token.balanceOf(l1Bridger.address);

  // Approve entire ERC-20 balance on L1
  const l1ApproveTx = await l1Token.approve(
    l1StandardBridge.address,
    l1BalanceStart,
  );
  await l1ApproveTx.wait();

  // Deposit entire token balance from L1 to L2
  const depositTx = await l1StandardBridge.bridgeERC20(
    l1Token.address,
    l2Token.address,
    l1BalanceStart,
    200_000,
    [],
  );

  const depositTxWithLogs = await depositTx.wait();
  console.log(depositTxWithLogs);
  const l1BalanceEnd = await l1Token.balanceOf(l1Bridger.address);

  const depositEvent = l1Portal.interface.parseLog(depositTxWithLogs.logs[3]);
  const depositMessage = {
    version: 0,
    nonce: depositEvent.args.nonce,
    sender: depositEvent.args.sender,
    target: depositEvent.args.target,
    value: depositEvent.args.value,
    gasLimit: depositEvent.args.gasLimit,
    data: depositEvent.args.data,
  };

  let blockNumber = await l1Provider.getBlockNumber();
  let rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);

  console.log("Initial block", { blockNumber, stateRoot, depositMessage });
  await waitUntilOracleBlock(l1Oracle, blockNumber);

  console.log({ depositHash: depositEvent.args.depositHash });
  let initiated = await l1Portal.initiatedDeposits(
    depositEvent.args.depositHash,
  );
  console.log({ initiated });

  const depositProof = await getDepositProof(
    l1Portal.address,
    depositEvent.args.depositHash,
    hexlifyBlockNum(blockNumber),
  );

  try {
    const finalizeTx = await l2Portal.finalizeDepositTransaction(
      blockNumber,
      depositMessage,
      depositProof.accountProof,
      depositProof.storageProof,
    );
    await finalizeTx.wait();
  } catch (e) {
    console.log({ e });
  }

  const l2BalanceEnd = await l2Token.balanceOf(l2Bridger.address);
  if (!l1BalanceEnd.eq(0) || !l2BalanceEnd.eq(l1BalanceStart)) {
    throw "unexpected end balance";
  }
  console.log("\tdeposited token...");

  // Approve entire ERC-20 balance on L2
  const l2ApproveTx = await l2Token.approve(
    l2StandardBridge.address,
    l2BalanceEnd,
  );
  await l2ApproveTx.wait();

  const withdrawalTx = await l2StandardBridge.bridgeERC20(
    l2Token.address,
    l1Token.address,
    l2BalanceEnd,
    200_000,
    [],
  );
  const txWithLogs = await withdrawalTx.wait();
  const withdrawBlockNum = txWithLogs.blockNumber;

  const l2BalanceEmpty = await l2Token.balanceOf(l2Bridger.address);
  if (!l2BalanceEmpty.eq(0)) {
    throw "unexpected L2 balance";
  }

  const withdrawEvent = l2Portal.interface.parseLog(txWithLogs.logs[3]);
  const withdrawMessage = {
    version: 0,
    nonce: withdrawEvent.args.nonce,
    sender: withdrawEvent.args.sender,
    target: withdrawEvent.args.target,
    value: withdrawEvent.args.value,
    gasLimit: withdrawEvent.args.gasLimit,
    data: withdrawEvent.args.data,
  };

  const withdrawalHash = withdrawEvent.args.withdrawalHash;
  console.log({ withdrawHash: withdrawalHash });
  initiated = await l2Portal.initiatedWithdrawals(withdrawalHash);
  console.log({ initiated });

  const [assertionId, assertionBlockNum] = await waitUntilBlockConfirmed(
    rollup,
    withdrawBlockNum,
  );

  // Get withdraw proof for the block the assertion committed to.
  const withdrawProof = await getWithdrawalProof(
    l2Portal.address,
    withdrawalHash,
    hexlifyBlockNum(assertionBlockNum),
  );

  // Get block for the block the assertion committed to.
  rawBlock = await l2Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(assertionBlockNum),
    false, // We only want the block header
  ]);
  let l2BlockHash = l2Provider.formatter.hash(rawBlock.hash);
  let l2StateRoot = l2Provider.formatter.hash(rawBlock.stateRoot);

  // Finalize withdraw
  console.log({ l2BlockHash, l2StateRoot });
  try {
    let finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
      withdrawMessage,
      assertionId,
      l2BlockHash,
      l2StateRoot,
      withdrawProof.accountProof,
      withdrawProof.storageProof,
    );
    console.log(finalizeTx);
    await finalizeTx.wait();
  } catch (e) {
    console.log({ e });
  }

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
