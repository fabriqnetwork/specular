import { getSignersAndContracts, deployTokenPair } from "../e2e/utils";

async function main() {
  const { l1Bridger, l2Relayer } = await getSignersAndContracts();

  const { l1Token, l2Token } = await deployTokenPair(l1Bridger, l2Relayer);
  console.log("\tdeployed tokens...");
  console.log("\tl1Token tokens...", l1Token.address);
  console.log("\tl2Token tokens...", l2Token.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
