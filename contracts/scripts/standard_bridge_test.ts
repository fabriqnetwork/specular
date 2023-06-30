import dotenv from "dotenv";
import { ethers } from "hardhat";

dotenv.config({ path: __dirname + "/../.env" });

function getStorageKey(messageHash: string) {
  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      ["bytes32", "uint256"],
      [messageHash, 0]
    )
  );
}

const l2Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:4011"
);

const l1Provider = new ethers.providers.JsonRpcProvider(
  "http://localhost:8545"
);

const l1BridgeAddr = "0xE7C2a73131dd48D8AC46dCD7Ab80C8cbeE5b410A";
const l2BridgeAddr = "0xF6168876932289D073567f347121A267095f3DD6";
const l1OracleAddress = "0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1";

async function main() {
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
    l1Bridger
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

  l1StandardBridge.on("ETHBridgeInitiated", (from, to, amount) => {
    console.log({
      event: "ETHBridgeInitiated",
      from,
      to,
      amount: ethers.utils.formatEther(amount),
    });
  });

  l2StandardBridge.on("ETHBridgeFinalized", (from, to, amount) => {
    console.log({
      event: "ETHBridgeFinalized",
      from,
      to,
      amount: ethers.utils.formatEther(amount),
    });
  });

  l2Portal.on("DepositFinalized", (hash, success) => {
    console.log({ event: "DepositFinalized", hash, success });
  });

  l1Oracle.on("L1OracleValuesUpdated", (number, hash) => {
    console.log({ event: "L1OracleValuesUpdated", number, hash });
  });

  const balanceStart = await l2Bridger.getBalance();
  const bridgeValue = ethers.utils.parseEther("0.1");

  const tx = await l1StandardBridge.bridgeETH(200_000, [], {
    value: bridgeValue,
  });
  await tx.wait();

  await new Promise((resolve, reject) => {
    l1Portal.once(
      "DepositInitiated",
      async (nonce, sender, target, value, gasLimit, data, depositHash) => {
        const crossDomainMessage = {
          version: 1,
          nonce,
          sender,
          target,
          value,
          gasLimit,
          data,
        };

        console.log({ event: "DepositInitiated", crossDomainMessage });

        const blockNumber = await l1Provider.getBlockNumber();
        const rawBlock = await l1Provider.send("eth_getBlockByNumber", [
          ethers.utils.hexValue(blockNumber),
          false, // We only want the block header
        ]);
        const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);
        await l1Oracle.setL1OracleValues(blockNumber, stateRoot, 0);

        const proof = await l1Provider.send("eth_getProof", [
          l1PortalAddress,
          [getStorageKey(depositHash)],
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
        resolve();
      }
    );
  });

  const balanceEnd = await l2Bridger.getBalance();
  if (!balanceEnd.sub(balanceStart).eq(bridgeValue)) {
    throw "unexpected end balance";
  }

  console.log("bridging ETH was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
