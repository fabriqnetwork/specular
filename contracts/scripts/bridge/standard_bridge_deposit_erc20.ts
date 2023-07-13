import { ethers } from "hardhat";
import { getSignersAndContracts, getStorageKey } from "./utils";

async function main() {
  const {
    l1Provider,
    l1Bridger,
    l2Bridger,
    l2Relayer,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2StandardBridge,
    l1Oracle,
  } = await getSignersAndContracts();

  const TestTokenFactory = await ethers.getContractFactory(
    "TestToken",
    l1Bridger
  );
  const l1Token = await TestTokenFactory.deploy();
  const l1BalanceStart = await l1Token.balanceOf(l1Bridger.address);

  const MintableERC20FactoryFactory = await ethers.getContractFactory(
    "OptimismMintableERC20Factory",
    l2Relayer
  );
  const mintableERC20Factory = await MintableERC20FactoryFactory.deploy(
    l2StandardBridge.address
  );
  const deployTx = await mintableERC20Factory.createOptimismMintableERC20(
    l1Token.address,
    "TestToken",
    "TT"
  );
  const deployTxWithLogs = await deployTx.wait();
  const deployEvent = await mintableERC20Factory.interface.parseLog(
    deployTxWithLogs.logs[0]
  );
  const l2TokenAddr = deployEvent.args.localToken;

  const MintableERC20Factory = await ethers.getContractFactory(
    "OptimismMintableERC20",
    l2Relayer
  );
  const l2Token = MintableERC20Factory.attach(l2TokenAddr);

  const approveTx = await l1Token.approve(
    l1StandardBridge.address,
    l1BalanceStart
  );
  await approveTx.wait();

  const bridgeTx = await l1StandardBridge.bridgeERC20(
    l1Token.address,
    l2Token.address,
    l1BalanceStart,
    200_000,
    []
  );
  const txWithLogs = await bridgeTx.wait();
  const l1BalanceEnd = await l1Token.balanceOf(l1Bridger.address);

  const initEvent = await l1Portal.interface.parseLog(txWithLogs.logs[3]);
  const crossDomainMessage = {
    version: 1,
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

  const proof = await l1Provider.send("eth_getProof", [
    l1Portal.address,
    [getStorageKey(initEvent.args.depositHash)],
    "latest",
  ]);
  const accountProof = proof.accountProof;
  const storageProof = proof.storageProof[0].proof;

  const tx = await l2Portal.finalizeDepositTransaction(
    crossDomainMessage,
    accountProof,
    storageProof
  );

  await tx.wait();

  const l2BalanceEnd = await l2Token.balanceOf(l2Bridger.address);

  if (!l1BalanceEnd.eq(0) || !l2BalanceEnd.eq(l1BalanceStart)) {
    throw "unexpected end balance";
  }

  console.log("bridging ERC20 was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
