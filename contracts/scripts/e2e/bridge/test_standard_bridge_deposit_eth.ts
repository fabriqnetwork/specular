import { ethers } from "hardhat";
import { BigNumber } from "ethers";
import { formatEther } from "ethers/lib/utils";
import {
  getSignersAndContracts,
  getDepositProof,
  delay,
} from "../utils";

async function main() {
  const {
    l1Provider,
    l2Provider,
    l2Bridger,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2StandardBridge,
    l1Oracle,
  } = await getSignersAndContracts();

  // TODO: portal should be funded as part of pre-deploy pipeline
  const donateTx = await l2Portal.donateETH({ value: ethers.utils.parseEther("1") })
  await donateTx;

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

  console.log("Initial block", { blockNumber, stateRoot, initEvent });

  let oracleStateRoot = await l1Oracle.stateRoot()
  while (oracleStateRoot !== stateRoot) {
    await delay(500)
    oracleStateRoot = await l1Oracle.stateRoot()
    console.log({ stateRoot, oracleStateRoot })
  }

  console.log({ depositHash: initEvent.args.depositHash })
  const initiated = await l1Portal.initiatedDeposits(initEvent.args.depositHash)
  console.log({ initiated })

  const onChainL1PortalAddr = await l2Portal.l1PortalAddress();
  console.log({ onChainL1PortalAddr, actualAddr: l1Portal.address })

  console.log({ L2BrideAddr: l2StandardBridge.address })

  const l2OtherBridge = await l2StandardBridge.OTHER_BRIDGE()
  const l2PortalAddr = await l2StandardBridge.PORTAL_ADDRESS()
  console.log({ l2OtherBridge, l1Bridge: l1StandardBridge.address, l2PortalAddr, l2PortalAddrActual: l2Portal.address })

  const l2PortalBalance = await l2Provider.getBalance(l2PortalAddr)
  console.log({ l2PortalBalance })

  const { accountProof, storageProof } = await getDepositProof(
    l1Portal.address,
    initEvent.args.depositHash,
    ethers.utils.hexlify(blockNumber)
  );

  try {
    const finalizeTx = await l2Portal.finalizeDepositTransaction(
      crossDomainMessage,
      accountProof,
      storageProof
    );
    await finalizeTx.wait();
  } catch(e) {
    console.log({ e })
  }

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
