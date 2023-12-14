// NOTE: this test will fail without portal level retryability
import { ethers } from "hardhat";
import { getSignersAndContracts, getDepositProof, hexlifyBlockNum } from "../utils";

async function main() {
  const {
    l1Provider,
    l2Provider,
    l2Bridger,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l1Oracle,
  } = await getSignersAndContracts();

  const balanceStart = await l2Bridger.getBalance();
  const portalBalance = await l2Provider.getBalance(l2Portal.address);
  const bridgeValue = portalBalance.add(10);

  console.log({ balanceStart: ethers.utils.formatEther(balanceStart) });
  console.log({ bridgeValue: ethers.utils.formatEther(bridgeValue) });

  const bridgeTx = await l1StandardBridge.bridgeETH(200_000, [], {
    value: bridgeValue,
  });
  const txWithLogs = await bridgeTx.wait();

  const initEvent = l1Portal.interface.parseLog(txWithLogs.logs[1]);
  const crossDomainMessage = {
    version: 0,
    nonce: initEvent.args.nonce,
    sender: initEvent.args.sender,
    target: initEvent.args.target,
    value: initEvent.args.value,
    gasLimit: initEvent.args.gasLimit,
    data: initEvent.args.data,
  };

  const blockNumber = await l1Provider.getBlockNumber();
  const rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);
  await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);

  const { accountProof, storageProof } = await getDepositProof(
    l1Portal.address,
    initEvent.args.depositHash,
    hexlifyBlockNum(blockNumber)
  );

  const finalizeTx = await l2Portal.finalizeDepositTransaction(
    crossDomainMessage,
    accountProof,
    storageProof
  );
  await finalizeTx.wait();

  const balanceEnd = await l2Bridger.getBalance();
  console.log({ balanceEnd: ethers.utils.formatEther(balanceEnd) });
  if (balanceEnd.sub(balanceStart).eq(bridgeValue)) {
    throw "unexpected end balance";
  }

  const donateTx = await l2Portal.donateETH({
    value: ethers.utils.parseUnits("1"),
  });
  await donateTx.wait();

  const retryTx = await l2Portal.finalizeDepositTransaction(
    crossDomainMessage,
    accountProof,
    storageProof
  );
  await retryTx.wait();

  const balanceFinal = await l2Bridger.getBalance();
  console.log({ balanceFinal: ethers.utils.formatEther(balanceFinal) });
  if (!balanceFinal.sub(balanceStart).eq(bridgeValue)) {
    throw "unexpected end balance";
  }

  console.log("bridging ETH was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
