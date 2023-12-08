const l2Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:4011"
);

const l1Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:8545"
);

const l1BridgeAddr = process.env.L1STANDARD_BRIDGE_ADDR;
const l2BridgeAddr = "0x2A00000000000000000000000000000000000012"
const l1OracleAddress = "0x2A00000000000000000000000000000000000010"
const rollupAddress = process.env.ROLLUP_ADDR;

export async function getSignersAndContracts() {
  const l1Bridger = new ethers.Wallet(
    "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
    l1Provider
  );
  const l2Bridger = new ethers.Wallet(
    "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
    l2Provider
  );

  const l1Relayer = new ethers.Wallet(
    "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
    l1Provider
  );
  const l2Relayer = new ethers.Wallet(
    "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
    l2Provider
  );

  const L1StandardBridgeFactory = await ethers.getContractFactory(
    "L1StandardBridge",
    l1Bridger
  );
  const l1StandardBridge = L1StandardBridgeFactory.attach(l1BridgeAddr);

  const L2StandardBridgeFactory = await ethers.getContractFactory(
    "L2StandardBridge",
    l2Bridger
  );
  const l2StandardBridge = L2StandardBridgeFactory.attach(l2BridgeAddr);

  const l1PortalAddress = await l1StandardBridge.PORTAL_ADDRESS();
  const L1PortalFactory = await ethers.getContractFactory(
    "L1Portal",
    l1Relayer
  );
  const l1Portal = L1PortalFactory.attach(l1PortalAddress);

  const l2PortalAddress = "0x2A00000000000000000000000000000000000011"; // await l2StandardBridge.PORTAL_ADDRESS();
  const L2PortalFactory = await ethers.getContractFactory(
    "L2Portal",
    l2Relayer
  );
  const l2Portal = L2PortalFactory.attach(l2PortalAddress);

  const L1OracleFactory = await ethers.getContractFactory(
    "L1Oracle",
    l2Relayer
  );
  const l1Oracle = L1OracleFactory.attach(l1OracleAddress);

  const RollupFactory = await ethers.getContractFactory("Rollup", l1Relayer);
  const rollup = await RollupFactory.attach(rollupAddress);

  const InboxFactory = await ethers.getContractFactory(
    "SequencerInbox",
    l1Relayer
  );
  const daProvider = await rollup.daProvider();
  const inbox = await InboxFactory.attach(daProvider);

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

export async function getDepositProof(portalAddress, depositHash, blockNumber="latest") {
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

export async function getWithdrawalProof(portalAddress, withdrawalHash) {
  const proof = await l2Provider.send("eth_getProof", [
    portalAddress,
    [getStorageKey(withdrawalHash)],
    "latest",
  ]);

  return {
    accountProof: proof.accountProof,
    storageProof: proof.storageProof[0].proof,
  };
}

export async function deployTokenPair(l1Bridger, l2Relayer) {
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
    l2BridgeAddr
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

  return {
    l1Token,
    l2Token,
  };
}

export function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function getLastBlockNumber(data) {
  // TODO: consider using the npm bindings package
  const InboxFactory = await ethers.getContractFactory("SequencerInbox");
  const iface = InboxFactory.interface;

  const decoded = iface.decodeFunctionData(
    data.slice(0, 10), // method id (8 hex chars) with leading Ox
    data
  );
  const contexts: BigNumber[] = decoded[0];
  const firstL2BlockNumber = decoded[2];
  const lastL2BlockNumber =
    contexts.length / 2 + firstL2BlockNumber.toNumber() - 1;
  return lastL2BlockNumber;
}

export function getStorageKey(messageHash: string) {
  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      ["bytes32", "uint256"],
      [messageHash, 0]
    )
  );
}
