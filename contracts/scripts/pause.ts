import { Contract, ethers, Wallet } from "ethers";
import { getSignersAndContracts } from "./e2e/utils";

async function main() {
  const { l1Provider, rollup, inbox, l1StandardBridge, l1Portal } =
    await getSignersAndContracts();

  const lastConfirmed = await rollup.lastConfirmedAssertionID();
  const lastCreated = await rollup.lastCreatedAssertionID();

  console.log({ lastConfirmed, lastCreated });
  return;

  if (!process.env.OWNER_PRIVATE_KEY) {
    throw new Error("no private key provided for contract owner");
  }
  const owner = new ethers.Wallet(process.env.OWNER_PRIVATE_KEY, l1Provider);

  // TODO: should this also pause L2 contracts?
  const contracts = [rollup, inbox, l1StandardBridge, l1Portal];

  if (process.env.UNPAUSE) {
    await unpauseContracts(contracts, owner);
    return;
  }

  await pauseContracts(contracts, owner);
}

async function pauseContracts(contracts: Contract[], owner: Wallet) {
  for (let contract of contracts) {
    const tx = await pauseContract(contract, owner);
    await tx.wait();
  }
}

async function unpauseContracts(contracts: Contract[], owner: Wallet) {
  for (let contract of contracts) {
    const tx = await unpauseContract(contract, owner);
    await tx.wait();
  }
}

async function pauseContract(contract: Contract, owner: Wallet) {
  const actualOwner = await contract.owner();
  const ownedContract = contract.connect(owner);
  console.log({
    contract: contract.address,
    actualOwner,
    providedOwner: owner.address,
  });
  return await ownedContract.pause();
}

async function unpauseContract(contract: Contract, owner: Wallet) {
  const actualOwner = await contract.owner();
  const ownedContract = contract.connect(owner);
  console.log({
    contract: contract.address,
    actualOwner,
    providedOwner: owner.address,
  });
  return await ownedContract.unpause();
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
