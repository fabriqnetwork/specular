import { ethers } from "hardhat";
import {
  getSignersAndContracts,
  getStorageKey,
  getDepositProof,
  getWithdrawalProof,
  delay,
  getLastBlockNumber,
  deployTokenPair,
} from "./utils";

async function main() {
  const { l1Bridger, l2Relayer } = await getSignersAndContracts();

  const { l1Token, l2Token } = await deployTokenPair(l1Bridger, l2Relayer);
  console.log("\tdeployed tokens...");
  console.log("\tl1Token tokens...", l1Token.address);
  console.log("\tl2Token tokens...", l2Token.address);
}
