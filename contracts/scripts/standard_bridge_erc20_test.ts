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

  l1StandardBridge.on(
    "ERC20BridgeInitiated",
    (localToken, remoteToken, from, to, amount) => {
      console.log({
        event: "ERC20BridgeInitiated",
        localToken,
        remoteToken,
        from,
        to,
        amount: ethers.utils.formatEther(amount),
      });
    }
  );

  l2StandardBridge.on(
    "ERC20BridgeFinalized",
    (localToken, remoteToken, from, to, amount) => {
      console.log({
        event: "ERC20BridgeFinalized",
        localToken,
        remoteToken,
        from,
        to,
        amount: ethers.utils.formatEther(amount),
      });
    }
  );

  l1Portal.on(
    "DepositInitiated",
    async (nonce, sender, target, value, gasLimit, data, depositHash) => {
      console.log({ event: "DepositInitiated", nonce, sender, target, value });
    }
  );

  l2Portal.on("DepositFinalized", (hash, success) => {
    console.log({ event: "DepositFinalized", hash, success });
  });

  l1Oracle.on("L1OracleValuesUpdated", (number, hash) => {
    console.log({ event: "L1OracleValuesUpdated", number, hash });
  });

  // deploy L1 ERC20
  const TestTokenFactory = await ethers.getContractFactory(
    "TestToken",
    l1Bridger
  );
  const l1Token = await TestTokenFactory.deploy();

  const l1BalanceStart = await l1Token.balanceOf(l1Bridger.address);
  console.log({ b: l1BalanceStart.toString() });

  // deploy L2 ERC20 Mintable Factory
  // create L2 ERC20 Mintable
  const MintableERC20FactoryFactory = await ethers.getContractFactory(
    "OptimismMintableERC20Factory",
    l2Relayer
  );
  const mintableERC20Factory = await MintableERC20FactoryFactory.deploy(
    l2BridgeAddr
  );
  await mintableERC20Factory.createOptimismMintableERC20(
    l1Token.address,
    "TestToken",
    "TT"
  );
  let l2TokenAddr = "";

  await new Promise((resolve, reject) => {
    mintableERC20Factory.once(
      "OptimismMintableERC20Created",
      (localToken, remoteToken) => {
        l2TokenAddr = localToken;
        console.log({ localToken, remoteToken });
        resolve();
      }
    );
  });

  const MintableERC20Factory = await ethers.getContractFactory(
    "OptimismMintableERC20",
    l2Relayer
  );
  const l2Token = MintableERC20Factory.attach(l2TokenAddr);

  const l2BalanceStart = await l2Token.balanceOf(l2Bridger.address);
  console.log({ b: l2BalanceStart.toString() });

  // TODO bridge from L1 -> L2
  await l1Token.approve(l1StandardBridge.address, l1BalanceStart);
  await l1StandardBridge.bridgeERC20(
    l1Token.address,
    l2Token.address,
    l1BalanceStart,
    200_000,
    []
  );

  const l1BalanceEnd = await l1Token.balanceOf(l1Bridger.address);
  console.log({ b: l1BalanceEnd.toString() });

  await new Promise((resolve, reject) => {
    l1Portal.on(
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

  const l2BalanceEnd = await l2Token.balanceOf(l2Bridger.address);
  console.log({ b: l2BalanceEnd.toString() });

  console.log("bridging ERC20 was successful");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
