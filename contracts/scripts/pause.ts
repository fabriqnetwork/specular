import { ethers } from "hardhat";
import * as addresses from "./e2e/addresses";

const rollupAddr = process.env.ROLLUP_ADDR;
const rollupPrKey = process.env.ROLLUP_PRIVATE_KEY;

const sequencerAddr = process.env.SEQUENCER_ADDR;
const sequencerPrKey = process.env.SEQUENCER_PRIVATE_KEY;

const l1BridgeAddr = process.env.L1STANDARD_BRIDGE;
const l1BridgePrKey = process.env.L1_BRIDGE_PRIVATE_KEY;

const l1OraclePrKey = process.env.L1_ORACLE_PRIVATE_KEY;

const l1PortalAddr = process.env.L1_PORTAL_ADDR;
const l1PortalPrKey = process.env.L1_PORTAL_PRIVATE_KEY;

const faucetAddr = process.env.FAUCET_ADDR;

const l1RPCAddr = process.env.L1_RPC_ADDRESS;
const l1Provider = new ethers.providers.JsonRpcProvider(l1RPCAddr);

const l2RPCAddr = process.env.L2_RPC_ADDRESS;
const l2Provider = new ethers.providers.JsonRpcProvider(l2RPCAddr);

async function pauseContract(
  contract: string,
  address: string,
  provider: ethers.providers.JsonRpcProvider,
  privateKey: string
) {
  const owner = new ethers.Wallet(privateKey, provider);
  const factory = await ethers.getContractFactory(contract, address);
  const c = factory.attach(address);
  const tx = c.pause();
  await tx.wait();
}

async function unpauseContract(
  contract: string,
  address: string,
  provider: ethers.providers.JsonRpcProvider,
  privateKey: string
) {
  const owner = new ethers.Wallet(privateKey, provider);
  const factory = await ethers.getContractFactory(contract, address);
  const c = factory.attach(address);
  const tx = c.unpause();
  await tx.wait();
}

async function pauseContracts() {
  // Pause L1 Contracts
  pauseContract("Rollup", rollupAddr, l1Provider, rollupPrKey);
  pauseContract("SequencerInbox", sequencerAddr, l1Provider, sequencerPrKey);
  pauseContract("L1StandardBridge", l1BridgeAddr, l1Provider, l1BridgePrKey);
  pauseContract("L1Oracle", addresses.l1OracleAddress, l1Provider, l1OraclePrKey);
  pauseContract("L1Portal", l1PortalAddr, l1Provider, l1PortalPrKey);

  // Pause L2 Contracts
  pauseContract("L2StandardBridge", addresses.l2StandardBridgeAddress, l2Provider, l2BridgePrKey);
  pauseContract("L2Portal", addresses.l2PortalAddress, l2Provider, l2PortalPrKey);
  pauseContract("Faucet", faucetAddr, l2Provider, faucetPrKey);
}

async function unpauseContracts() {
  // Unpause L1 Contracts
  unpauseContract("Rollup", rollupAddr, l1Provider, rollupPrKey);
  unpauseContract("SequencerInbox", sequencerAddr, l1Provider, sequencerPrKey);
  unpauseContract("L1StandardBridge", l1BridgeAddr, l1Provider, l1BridgePrKey);
  unpauseContract("L1Oracle", addresses.l1OracleAddress, l1Provider, l1OraclePrKey);
  unpauseContract("L1Portal", l1PortalAddr, l1Provider, l1PortalPrKey);

  // Unpause L2 Contracts
  unpauseContract("L2StandardBridge", addresses.l2StandardBridgeAddress, l2Provider, l2BridgePrKey);
  unpauseContract("L2Portal", addresses.l2PortalAddress, l2Provider, l2PortalPrKey);
  unpauseContract("Faucet", faucetAddr, l2Provider, faucetPrKey);
}

pauseContracts()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
