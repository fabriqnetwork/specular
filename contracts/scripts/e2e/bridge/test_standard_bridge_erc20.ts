import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getDepositProof,
  getWithdrawalProof,
  delay,
  deployTokenPair,
  hexlifyBlockNum,
} from "../utils";

import { BigNumber } from "ethers";

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
    inbox,
  } = await getSignersAndContracts();

  const { l1Token, l2Token } = await deployTokenPair(l1Bridger, l2Relayer);
  console.log(`Deployed tokens ${l1Token.address}, ${l2Token.address}`);

  // TODO: portal should be funded as part of pre-deploy pipeline
  const donateTx = await l2Portal.donateETH({ value: ethers.utils.parseEther("1") })
  await donateTx;

  const l1BalanceStart = await l1Token.balanceOf(l1Bridger.address);
  const l2BalanceStart = await l2Token.balanceOf(l2Relayer.address);

  console.log({L1Bridger: l1Bridger.address, l2Replayer: l2Relayer.address, L1BalanceStart: l1BalanceStart, l2BalanceStart: l2BalanceStart})

  const approveTx = await l1Token.approve(
    l1StandardBridge.address,
    l1BalanceStart
  );
  const approveTxWithLogs = await approveTx.wait();
  console.log({ApproveTx: approveTxWithLogs})

  const depositTx = await l1StandardBridge.bridgeERC20(
    l1Token.address,
    l2Token.address,
    l1BalanceStart,
    200_000,
    []
  );

  const depositTxWithLogs = await depositTx.wait();
  console.log(depositTxWithLogs)
  const l1BalanceEnd = await l1Token.balanceOf(l1Bridger.address);

  const depositEvent = l1Portal.interface.parseLog(
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

  let blockNumber = await l1Provider.getBlockNumber();
  const rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);

  console.log("Initial block", { blockNumber, stateRoot, depositMessage });

  let oracleStateRoot = await l1Oracle.stateRoot()
  while (oracleStateRoot !== stateRoot) {
    await delay(500)
    oracleStateRoot = await l1Oracle.stateRoot()
    console.log({ stateRoot, oracleStateRoot })
  }

  console.log({ depositHash: depositEvent.args.depositHash })
  const initiated = await l1Portal.initiatedDeposits(depositEvent.args.depositHash)
  console.log({ initiated })

  const onChainL1PortalAddr = await l2Portal.l1PortalAddress();
  console.log({ onChainL1PortalAddr, actualAddr: l1Portal.address })

  console.log({ L2BridgeAddr: l2StandardBridge.address })

  const l2OtherBridge = await l2StandardBridge.OTHER_BRIDGE()
  const l2PortalAddr = await l2StandardBridge.PORTAL_ADDRESS()
  console.log({ l2OtherBridge, l1Bridge: l1StandardBridge.address, l2PortalAddr, l2PortalAddrActual: l2Portal.address })

  const l2PortalBalance = await l2Provider.getBalance(l2PortalAddr)
  console.log({ l2PortalBalance })

  const depositProof = await getDepositProof(
    l1Portal.address,
    depositEvent.args.depositHash,
    hexlifyBlockNum(blockNumber),
  );

  try {
    const finalizeTx = await l2Portal.finalizeDepositTransaction(
      depositMessage,
      depositProof.accountProof,
      depositProof.storageProof
    );
    await finalizeTx.wait();
  } catch(e) {
    console.log({ e })
  }

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
  console.log(withdrawalTx)
  const txWithLogs = await withdrawalTx.wait();

  const l2BalanceEmpty = await l2Token.balanceOf(l2Bridger.address);
  if (!l2BalanceEmpty.eq(0)) {
    throw "unexpected L2 balance";
  }

  const initEvent = l2Portal.interface.parseLog(txWithLogs.logs[3]);
  const crossDomainMessage = {
    version: 0,
    nonce: initEvent.args.nonce,
    sender: initEvent.args.sender,
    target: initEvent.args.target,
    value: initEvent.args.value,
    gasLimit: initEvent.args.gasLimit,
    data: initEvent.args.data,
  };

  blockNumber = txWithLogs.blockNumber;

  let assertionId: number | undefined = undefined;
  let lastConfirmedBlockNum: number | undefined = undefined;

  inbox.on(
    inbox.filters.TxBatchAppended(),
    (event) => {
      console.log(`TxBatchAppended blockNum: ${event.blockNumber}`)
    }
  );

  rollup.on(rollup.filters.AssertionConfirmed(), async (id: BigNumber) => {
    console.log("AssertionConfirmed", "id", assertionId)
    if (!assertionId) {
      assertionId = id.toNumber();
      const assertion = await rollup.getAssertion(assertionId);
      lastConfirmedBlockNum = assertion.blockNum.toNumber();
    }
  });

  console.log("Waiting for assertion to be confirmed...");
  while (!assertionId || !lastConfirmedBlockNum || lastConfirmedBlockNum < blockNumber) {
    await delay(500);
  }

  const { accountProof, storageProof } = await getWithdrawalProof(
    l2Portal.address,
    initEvent.args.withdrawalHash
  );

  let l2VmHash = l2Provider.formatter.hash(rawBlock.hash);
  const finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
    crossDomainMessage,
    assertionId,
    l2VmHash,
    accountProof,
    storageProof
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
