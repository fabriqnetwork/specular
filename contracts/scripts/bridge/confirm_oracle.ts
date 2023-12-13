import { ethers } from "hardhat";
import { getSignersAndContracts } from "../e2e/utils";

async function main() {
  const { l1Provider, l1Oracle } = await getSignersAndContracts();

  const blockNumber = await l1Provider.getBlockNumber();
  const rawBlock = await l1Provider.send("eth_getBlockByNumber", [
    ethers.utils.hexValue(blockNumber),
    false, // We only want the block header
  ]);
  const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);
  await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
