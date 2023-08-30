const l2Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:4011"
);

const l1Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:8545"
);

const l1BridgeAddr = process.env.L1_BRIDGE_ADDR;
const l2BridgeAddr = process.env.L2_BRIDGE_ADDR;
const l1OracleAddress = process.env.L1_ORACLE_ADDR;
const rollupAddress = process.env.ROLLUP_ADDR;

export async function getSignersAndContracts() {
  const l1Bridger = new ethers.Wallet(
    process.env.BRIDGER_PRIVATE_KEY,
    l1Provider
  );
  const l2Bridger = new ethers.Wallet(
    process.env.BRIDGER_PRIVATE_KEY,
    l2Provider
  );

  const l1Relayer = new ethers.Wallet(
    process.env.SEQUENCER_PRIVATE_KEY,
    l1Provider
  );
  const l2Relayer = new ethers.Wallet(
    process.env.SEQUENCER_PRIVATE_KEY,
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

  const l2PortalAddress = await l2StandardBridge.PORTAL_ADDRESS();
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

export function getStorageKey(messageHash: string) {
  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      ["bytes32", "uint256"],
      [messageHash, 0]
    )
  );
}
