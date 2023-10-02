import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getStorageKey,
  getDepositProof,
  delay,
} from "./utils";

async function main() {
  const { l2Relayer, l2Bridger } = await getSignersAndContracts();

  const value = ethers.utils.parseEther("0.1");

  const startBalances = {
    l2Relayer: await l2Relayer.getBalance(),
    l2Bridger: await l2Bridger.getBalance(),
  };

  // TODO: should we randomize numTx and value?
  const numTx = 5;
  for (let i = 0; i < numTx; i++) {
    const tx = await l2Relayer.sendTransaction({
      to: l2Bridger.address,
      value,
    });
    await tx.wait();
  }

  const endBalances = {
    l2Relayer: await l2Relayer.getBalance(),
    l2Bridger: await l2Bridger.getBalance(),
  };

  console.log({ startBalances, endBalances });

  const totalValue = value.mul(numTx);

  if (!endBalances.l2Bridger.sub(startBalances.l2Bridger).eq(totalValue)) {
    throw "unexpected end balance";
  }

  const error = ethers.utils.parseEther("0.001");
  if (
    !startBalances.l2Relayer
      .sub(endBalances.l2Relayer)
      .sub(totalValue)
      .lt(error)
  ) {
    throw "unexpected end balance";
  }

  console.log("transactions test was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
