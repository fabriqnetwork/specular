import { ethers } from "hardhat";
import { BigNumber } from "ethers";
import { formatEther } from "ethers/lib/utils";
import {
  getSignersAndContracts,
  getDepositProof,
  hexlifyBlockNum,
  waitUntilOracleBlock,
} from "../utils";

async function main() {
  const {
    l1Provider,
    l2Bridger,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l1Oracle,
  } = await getSignersAndContracts();

  // TODO: portal should be funded as part of pre-deploy pipeline
  const donateTx = await l2Portal.donateETH({
    value: ethers.utils.parseEther("1"),
  });
  await donateTx;

  const balanceStart: BigNumber = await l2Bridger.getBalance();
  const bridgeValue: BigNumber = ethers.utils.parseEther("0.1");

  const bridgeTx = await l1StandardBridge.bridgeETH(200_000, [], {
    value: bridgeValue,
  });
  const txWithLogs = await bridgeTx.wait();

  const depositEvent = l1Portal.interface.parseLog(txWithLogs.logs[1]);
  const despositMessage = {
    version: 0,
    nonce: depositEvent.args.nonce,
    sender: depositEvent.args.sender,
    target: depositEvent.args.target,
    value: depositEvent.args.value,
    gasLimit: depositEvent.args.gasLimit,
    data: depositEvent.args.data,
  };

  let depositBlockNumber = await l1Provider.getBlockNumber();
  await waitUntilOracleBlock(l1Oracle, depositBlockNumber);

  console.log({ depositHash: depositEvent.args.depositHash });
  const initiated = await l1Portal.initiatedDeposits(
    depositEvent.args.depositHash,
  );
  console.log({ initiated });

  const depositProof = await getDepositProof(
    l1Portal.address,
    depositEvent.args.depositHash,
    hexlifyBlockNum(depositBlockNumber),
  );

  try {
    const finalizeTx = await l2Portal.finalizeDepositTransaction(
      depositBlockNumber,
      despositMessage,
      depositProof.accountProof,
      depositProof.storageProof,
    );
    await finalizeTx.wait();
  } catch (e) {
    console.log({ e });
  }

  // balanceStart <- balance on L2 before bridging
  // balanceEnd <- balance on L2 after bridging
  // bridgeValue <- the value we're transfering
  // Expected: balanceEnd - balanceStart - bridgeValue ~== 0
  const balanceEnd: BigNumber = await l2Bridger.getBalance();
  const actualDiff: BigNumber = balanceEnd
    .sub(balanceStart)
    .sub(bridgeValue)
    .abs();
  const acceptableMargin: BigNumber = ethers.utils.parseEther("0.0001");

  if (!actualDiff.lt(acceptableMargin)) {
    const situation = {
      balanceStart: formatEther(balanceStart),
      bridgeValue: formatEther(bridgeValue),
      balanceEnd: formatEther(balanceEnd),
      actualDiff: formatEther(actualDiff),
      acceptableMargin: formatEther(acceptableMargin),
    };
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
