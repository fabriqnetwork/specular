import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getWithdrawalProof,
  delay,
  hexlifyBlockNum,
} from "../utils";
import { BigNumber } from "ethers";

async function main() {
  const {
    l1Provider,
    l1Bridger,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2StandardBridge,
    l2Provider,
    rollup,
    inbox,
  } = await getSignersAndContracts();

  const donateTx = await l1Portal.donateETH({
    value: ethers.utils.parseEther("1"),
  });
  await donateTx.wait();

  const donateTx2 = await l2Portal.donateETH({
    value: ethers.utils.parseEther("1"),
  });
  await donateTx2.wait();

  const balanceStart = await l1Bridger.getBalance();
  const bridgeValue = ethers.utils.parseEther("0.1");

  const bridgeTx = await l2StandardBridge.bridgeETH(200_000, [], {
    value: bridgeValue,
  });
  const txWithLogs = await bridgeTx.wait();
  const blockNumber = txWithLogs.blockNumber;

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

  const withdrawalHash = initEvent.args.withdrawalHash
  console.log({ withdrawHash: withdrawalHash })
  const initiated = await l2Portal.initiatedWithdrawals(withdrawalHash)
  console.log({ initiated })

  console.log({ L2BridgeAddr: l2StandardBridge.address })

  const l1OtherBridge = await l2StandardBridge.OTHER_BRIDGE()
  const l1PortalAddr = await l1StandardBridge.PORTAL_ADDRESS()
  console.log({ l1OtherBridge, l1Bridge: l1StandardBridge.address, l1PortalAddr, l2PortalAddrActual: l1Portal.address })

  const l1PortalBalance = await l1Provider.getBalance(l1PortalAddr)
  console.log({ l1PortalBalance })

  let assertionId: number | undefined = undefined;
  let lastConfirmedBlockNum: number | undefined = undefined;

  inbox.on(
    inbox.filters.TxBatchAppended(),
    (event) => {
      console.log(`TxBatchAppended blockNum: ${event.blockNumber}`)
    }
  );

  rollup.on(rollup.filters.AssertionConfirmed(), async (id: BigNumber) => {
    console.log("AssertionConfirmed")
    if (!assertionId) {
      assertionId = id.toNumber();
      const assertion = await rollup.getAssertion(id.toNumber());
      lastConfirmedBlockNum = assertion.blockNum.toNumber();
    }
  });

  rollup.on(rollup.filters.AssertionCreated(), () => {
    console.log("AssertionCreated")
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
 
  console.log(`Assertion confirmed: ${assertionId}, l2 block: ${lastConfirmedBlockNum}`);

  const { accountProof, storageProof } = await getWithdrawalProof(
    l2Portal.address,
    withdrawalHash,
    hexlifyBlockNum(lastConfirmedBlockNum)
  );

  let rawBlock = await l2Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(lastConfirmedBlockNum),
    false, // We only want the block header
  ]);

  let l2StateRoot = l2Provider.formatter.hash(rawBlock.stateRoot);
  let l2BlockHash = l2Provider.formatter.hash(rawBlock.hash);
  console.log("Finalizing withdraw", "stateRoot", stateRoot, "blockNum", lastConfirmedBlockNum);
  try {
    let finalizeTx = await l1Portal.finalizeWithdrawalTransaction(
      crossDomainMessage,
      assertionId,
      l2StateRoot,
      l2BlockHash,
      accountProof,
      storageProof
    );
    console.log(finalizeTx)
    await finalizeTx.wait();
  } catch(e) {
    console.log({ e })
  }

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
