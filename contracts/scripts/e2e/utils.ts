import * as addresses from "./addresses";
import { Wallet, BigNumber, Contract } from "ethers";
import { ethers } from "hardhat";

const l1Provider = new ethers.providers.JsonRpcProvider(
  process.env.L2_RPC_URL || "http://localhost:8545",
);

const l2Provider = new ethers.providers.JsonRpcProvider(
  process.env.L1_RPC_URL || "http://localhost:4011",
);

const l1BridgeAddress = process.env.L1STANDARD_BRIDGE_ADDR || "invalid address";
const rollupAddress = process.env.ROLLUP_ADDR || "invalid address";

export async function getSignersAndContracts() {
  const l1Bridger = new ethers.Wallet(
    process.env.BRIDGER_PRIVATE_KEY || "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
    l1Provider,
  );
  const l2Bridger = new ethers.Wallet(
    process.env.BRIDGER_PRIVATE_KEY || "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
    l2Provider,
  );

  const l1Relayer = new ethers.Wallet(
    process.env.BRIDGER_PRIVATE_KEY || "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
    l1Provider,
  );
  const l2Relayer = new ethers.Wallet(
    process.env.BRIDGER_PRIVATE_KEY || "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
    l2Provider,
  );

  const L1StandardBridgeFactory = await ethers.getContractFactory(
    "L1StandardBridge",
    l1Bridger,
  );
  const l1StandardBridge = L1StandardBridgeFactory.attach(l1BridgeAddress);

  const L2StandardBridgeFactory = await ethers.getContractFactory(
    "L2StandardBridge",
    l2Bridger,
  );
  const l2StandardBridge = L2StandardBridgeFactory.attach(
    addresses.l2StandardBridgeAddress,
  );

  const l1PortalAddress = await l1StandardBridge.PORTAL_ADDRESS();
  const L1PortalFactory = await ethers.getContractFactory(
    "L1Portal",
    l1Relayer,
  );
  const l1Portal = L1PortalFactory.attach(l1PortalAddress);

  const L2PortalFactory = await ethers.getContractFactory(
    "L2Portal",
    l2Relayer,
  );
  const l2Portal = L2PortalFactory.attach(addresses.l2PortalAddress);

  const L1OracleFactory = await ethers.getContractFactory(
    "L1Oracle",
    l2Relayer,
  );
  const l1Oracle = L1OracleFactory.attach(addresses.l1OracleAddress);

  const RollupFactory = await ethers.getContractFactory("Rollup", l1Relayer);
  const rollup = RollupFactory.attach(rollupAddress);

  const InboxFactory = await ethers.getContractFactory(
    "SequencerInbox",
    l1Relayer,
  );
  const daProvider = await rollup.daProvider();
  const inbox = InboxFactory.attach(daProvider);

  // l1Portal.on("*", (...args) => console.log({ ...args }));
  // l2Portal.on("*", (...args) => console.log({ ...args }));
  // l1StandardBridge.on("*", (...args) => console.log({ ...args }));
  // l2StandardBridge.on("*", (...args) => console.log({ ...args }));

  return {
    l1Provider,
    l2Provider,
    l1Bridger,
    l2Bridger,
    l1Relayer,
    l2Relayer,
    l1Portal,
    l2Portal,
    l1StandardBridge,
    l2StandardBridge,
    l1Oracle,
    rollup,
    inbox,
  };
}

export async function getDepositProof(
  portalAddress: string,
  depositHash: string,
  blockNumber: string,
) {
  const proof = await l1Provider.send("eth_getProof", [
    portalAddress,
    [getStorageKey(depositHash)],
    blockNumber,
  ]);

  return {
    accountProof: proof.accountProof,
    storageProof: proof.storageProof[0].proof,
  };
}

export async function getWithdrawalProof(
  portalAddress: string,
  withdrawalHash: string,
  blockNumber: string,
) {
  const proof = await l2Provider.send("eth_getProof", [
    portalAddress,
    [getStorageKey(withdrawalHash)],
    blockNumber,
  ]);

  return {
    accountProof: proof.accountProof,
    storageProof: proof.storageProof[0].proof,
  };
}

export async function deployTokenPair(l1Bridger: Wallet, l2Relayer: Wallet) {
  const TestTokenFactory = await ethers.getContractFactory(
    "TestToken",
    l1Bridger,
  );
  const l1Token = await TestTokenFactory.deploy();

  const MintableERC20FactoryFactory = await ethers.getContractFactory(
    "MintableERC20Factory",
    l2Relayer,
  );
  const mintableERC20Factory = await MintableERC20FactoryFactory.deploy(
    addresses.l2StandardBridgeAddress,
  );
  const deployTx = await mintableERC20Factory.createMintableERC20(
    l1Token.address,
    "TestToken",
    "TT",
  );
  const deployTxWithLogs = await deployTx.wait();
  console.log(deployTxWithLogs);
  const deployEvent = mintableERC20Factory.interface.parseLog(
    deployTxWithLogs.logs[0],
  );
  const l2TokenAddr = deployEvent.args.localToken;

  const MintableERC20Factory = await ethers.getContractFactory(
    "MintableERC20",
    l2Relayer,
  );
  const l2Token = MintableERC20Factory.attach(l2TokenAddr);

  return {
    l1Token,
    l2Token,
  };
}

export function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function getStorageKey(messageHash: string) {
  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      ["bytes32", "uint256"],
      [messageHash, 0],
    ),
  );
}

/**
 * Blocking function that only exits once the block relayed by L1Oracle is >= the blockNumber
 * @param {Contract} l1Oracle - the oracle contract deployed on L2
 * @param {number} blockNumber - the block we are waiting for
 */
export async function waitUntilOracleBlock(
  l1Oracle: Contract,
  blockNumber: number,
) {
  console.log(`Waiting for L1Oracle to relay block #${blockNumber}...`);
  let oracleBlockNumber = await l1Oracle.number();
  while (oracleBlockNumber < blockNumber) {
    await delay(500);
    oracleBlockNumber = await l1Oracle.number();
  }
}

/**
 * Blocking function that only exits once L2 block with blockNum has been confirmed on L1
 * @param {Contract} rollup - the rollup contract deployed on L1
 * @param {number} blockNum - number of the L2 block we are waiting on
 */
export async function waitUntilBlockConfirmed(
  rollup: Contract,
  blockNum: number,
): Promise<[number, number]> {
  let confirmedAssertionId: number | undefined = undefined;
  let confirmedBlockNum: number;

  rollup.on(rollup.filters.AssertionCreated(), () => {
    console.log("AssertionCreated");
  });

  rollup.on(rollup.filters.AssertionConfirmed(), async (id: BigNumber) => {
    const assertionId = id.toNumber();
    const assertion = await rollup.getAssertion(assertionId);
    const assertionBlockNum = assertion.blockNum.toNumber();

    console.log({
      msg: "AssertionConfirmed",
      id: assertionId,
      blockNum: assertionBlockNum,
    });
    if (!confirmedAssertionId && blockNum <= assertionBlockNum) {
      // Found the first assertion to confirm block
      confirmedAssertionId = assertionId;
      confirmedBlockNum = assertionBlockNum;
    }
  });

  console.log(`Waiting for L2 block ${blockNum} to be confirmed...`);
  while (!confirmedAssertionId) {
    await delay(500);
  }

  return [confirmedAssertionId!, confirmedBlockNum!];
}

// eth_getProof block number param cannot have leading zeros
// This function hexlifys blockNum and strips leading zeros
export function hexlifyBlockNum(blockNum: number): string {
  let hexBlockNum = ethers.utils.hexlify(blockNum);
  // Check if the string starts with "0x" and contains more than just "0x".
  if (hexBlockNum.startsWith("0x") && hexBlockNum.length > 2) {
    let strippedString = "0x";

    // Iterate through the characters of the input string starting from the third character (index 2).
    for (let i = 2; i < hexBlockNum.length; i++) {
      if (hexBlockNum[i] !== "0") {
        strippedString += hexBlockNum.substring(i); // Append the remaining characters.
        return strippedString;
      }
    }

    // If all characters are '0', return "0x0".
    return "0x0";
  }

  // If the input is not in the expected format, return it as is.
  return hexBlockNum;
}
