import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getWithdrawalProof,
  hexlifyBlockNum,
  waitUntilBlockConfirmed,
} from "../utils";
async function main() {
  const {
    l1Bridger,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2StandardBridge,
    l2Provider,
    rollup,
  } = await getSignersAndContracts();

  // TODO: portal should be funded as part of deploy pipeline
  const donateTx = await l1Portal.donateETH({
    value: ethers.utils.parseEther("1"),
  });
  await donateTx.wait();

  l1StandardBridge.on(
    l1StandardBridge.filters.ETHBridgeFinalized(),
    async (from, to, amount, data) => {
      console.log({ msg: "ETHBridgeFinalized", from, to, amount, data });
    }
  );

  const balanceStart = await l1Bridger.getBalance();
  const bridgeValue = ethers.utils.parseEther("0.1");

  // Bridge ETH from L2 to L1
  const bridgeTx = await l2StandardBridge.bridgeETH(200_000, [], {
    value: bridgeValue,
  });
  const txWithLogs = await bridgeTx.wait();
  const withdrawTxBlockNum = txWithLogs.blockNumber;

  const withdrawEvent = l2Portal.interface.parseLog(txWithLogs.logs[1]);
  const withdrawMessage = {
    version: 0,
    nonce: withdrawEvent.args.nonce,
    sender: withdrawEvent.args.sender,
    target: withdrawEvent.args.target,
    value: withdrawEvent.args.value,
    gasLimit: withdrawEvent.args.gasLimit,
    data: withdrawEvent.args.data,
  };

  const withdrawalHash = withdrawEvent.args.withdrawalHash
  console.log({ withdrawHash: withdrawalHash })
  const initiated = await l2Portal.initiatedWithdrawals(withdrawalHash)
  console.log({ initiated })

  const [assertionId, assertionBlockNum] = await waitUntilBlockConfirmed(rollup, withdrawTxBlockNum)


  // Get withdraw proof for the block the assertion committed to.
  const withdrawProof = await getWithdrawalProof(
    l2Portal.address,
    withdrawalHash,
    hexlifyBlockNum(assertionBlockNum)
  );

  // Get block for the block the assertion committed to.
  let rawBlock = await l2Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(assertionBlockNum),
    false, // We only want the block header
  ]);
  let l2BlockHash = l2Provider.formatter.hash(rawBlock.hash);
  let l2StateRoot = l2Provider.formatter.hash(rawBlock.stateRoot);

  // Finalize withdraw
  console.log({l2BlockHash, l2StateRoot});
  try {
    let finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
      withdrawMessage,
      assertionId,
      l2BlockHash,
      l2StateRoot,
      withdrawProof.accountProof,
      withdrawProof.storageProof
    );
    console.log(finalizeTx)
    await finalizeTx.wait();
  } catch(e) {
    console.log({ e })
  }

  // Confirm ETH balance was bridged
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
