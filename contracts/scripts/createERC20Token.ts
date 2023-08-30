import { ethers } from "hardhat";
import { getSignersAndContracts, getStorageKey } from "./utils";

async function main() {
  const { l1Bridger, l2Relayer, l2StandardBridge } =
    await getSignersAndContracts();

  const TestTokenFactory = await ethers.getContractFactory(
    "TestToken",
    l1Bridger
  );
  const l1Token = await TestTokenFactory.deploy();

  const MintableERC20FactoryFactory = await ethers.getContractFactory(
    "MintableERC20Factory",
    l2Relayer
  );
  const mintableERC20Factory = await MintableERC20FactoryFactory.deploy(
    l2StandardBridge.address
  );
  const deployTx = await mintableERC20Factory.createMintableERC20(
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
    "MintableERC20",
    l2Relayer
  );
  const l2Token = MintableERC20Factory.attach(l2TokenAddr);

  console.log("L1 Token Addresss:", l1Token.address);
  console.log("L2 Token Addresss:", l2Token.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
