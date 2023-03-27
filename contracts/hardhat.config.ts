import "dotenv/config";
import "@nomiclabs/hardhat-ethers";
import { ethers } from "ethers";
import { HardhatUserConfig } from "hardhat/types";
import "@nomiclabs/hardhat-etherscan";
import "@nomiclabs/hardhat-waffle";
import "hardhat-abi-exporter";
import "hardhat-deploy";
import "@openzeppelin/hardhat-upgrades";

const mnemonic =
  process.env.MNEMONIC ??
  "test test test test test test test test test test test junk";

const sequencerWallet = ethers.Wallet.fromMnemonic(mnemonic);
// path for second account from mnemonic
const deployerPath = "m/44'/60'/1'/0/0";
const deployerWallet = ethers.Wallet.fromMnemonic(mnemonic, deployerPath);
const INFURA_KEY = process.env.INFURA_KEY || "";
const ALCHEMY_KEY = process.env.ALCHEMY_KEY || "";
const SEQUENCER_PRIVATE_KEY = process.env.SEQUENCER_PRIVATE_KEY;
const DEPLOYER_PRIVATE_KEY = process.env.DEPLOYER_PRIVATE_KEY;
const ETHERSCAN_API_KEY = process.env.ETHERSCAN_API_KEY || "";
const DEPLOYER_ADDRESS = process.env.DEPLOYER_ADDRESS ?? deployerWallet.address;
const SEQUENCER_ADDRESS =
  process.env.SEQUENCER_ADDRESS ?? sequencerWallet.address;

function createConfig(network: string) {
  return {
    url: getNetworkURL(network),
    accounts: getNetworkAccounts(network),
  };
}

function getNetworkAccounts(network: string) {
  if (network === "goerli" || network === "chiado") {
    return SEQUENCER_PRIVATE_KEY && DEPLOYER_PRIVATE_KEY
      ? [`0x${SEQUENCER_PRIVATE_KEY}`, `0x${DEPLOYER_PRIVATE_KEY}`]
      : { mnemonic };
  } else {
    return { mnemonic };
  }
}

function getNetworkURL(network: string) {
  if (network === "mainnet") {
    return `https://mainnet.infura.io/v3/${INFURA_KEY}`;
  } else if (network === "goerli") {
    return `https://goerli.infura.io/v3/${INFURA_KEY}`;
  } else if (network === "sepolia") {
    return `https://sepolia.infura.io/v3/${INFURA_KEY}`;
  } else if (network === "gnosis") {
    return "https://rpc.gnosischain.com/";
  } else if (network === "chiado") {
    return "https://rpc.chiadochain.net";
  }
}

const config: HardhatUserConfig = {
  defaultNetwork: "hardhat",
  solidity: {
    version: "0.8.4",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  paths: {
    sources: "./src",
  },
  abiExporter: {
    path: "./abi",
    runOnCompile: true,
    clear: true,
  },
  namedAccounts: {
    sequencer: {
      default: 0,
      goerli: SEQUENCER_ADDRESS,
      chiado: SEQUENCER_ADDRESS,
    },
    validator: 1,
    deployer: {
      default: 2,
      goerli: DEPLOYER_ADDRESS,
      chiado: DEPLOYER_ADDRESS,
    },
  },
  networks: {
    mainnet: createConfig("mainnet"),
    sepolia: createConfig("sepolia"),
    goerli: createConfig("goerli"),
    gnosis: createConfig("gnosis"),
    chiado: createConfig("chiado"),
    hardhat: {
      gas: "auto",
      mining: {
        auto: true,
        interval: 5000,
        mempool: {
          order: "fifo",
        },
      },
    },
  },
};

export default config;
