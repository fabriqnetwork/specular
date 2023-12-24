import "dotenv/config";
import "@nomiclabs/hardhat-ethers";
import { ethers } from "ethers";
import { HardhatUserConfig, HttpNetworkUserConfig } from "hardhat/types";
import "@nomiclabs/hardhat-etherscan";
import "@nomiclabs/hardhat-waffle";
import "hardhat-abi-exporter";
import "hardhat-deploy";
import "@openzeppelin/hardhat-upgrades";

const mnemonic =
  process.env.MNEMONIC ??
  "test test test test test test test test test test test junk";
const wallet = ethers.Wallet.fromMnemonic(mnemonic);

const INFURA_KEY = process.env.INFURA_KEY || "";
const DEPLOYER_PRIVATE_KEY = process.env.DEPLOYER_PRIVATE_KEY || "";
const DEPLOYER_ADDRESS = process.env.DEPLOYER_ADDRESS || wallet.address;

function createConfig(baseConfig: HardhatUserConfig) {
  if (baseConfig.networks && DEPLOYER_PRIVATE_KEY) {
    baseConfig.networks["mainnet"] = createNetworkConfig("mainnet");
    baseConfig.networks["sepolia"] = createNetworkConfig("sepolia");
    baseConfig.networks["localhost"] = createNetworkConfig("localhost");
  } else {
    console.warn("DEPLOYER_PRIVATE_KEY not found, only exporting network `localhost`.")
  }
  return baseConfig;
}

function createNetworkConfig(network: string): HttpNetworkUserConfig {
  return {
    url: getNetworkURL(network),
    accounts: network === "hardhat" ? { mnemonic } : [DEPLOYER_PRIVATE_KEY],
    live: network !== "localhost" && network !== "hardhat",
    saveDeployments: true,
    deploy: ["deploy/l1"],
  };
}

function getNetworkURL(network: string) {
  if (network === "mainnet") {
    return `https://mainnet.infura.io/v3/${INFURA_KEY}`;
  } else if (network === "sepolia") {
    // https://rpc.sepolia.org/
    return `https://sepolia.infura.io/v3/${INFURA_KEY}`;
  }
  return "http://0.0.0.0:8545";
}

const baseConfig: HardhatUserConfig = {
  defaultNetwork: "hardhat",
  solidity: {
    version: "0.8.17",
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
    deployer: {
      default: DEPLOYER_ADDRESS,
      localhost: DEPLOYER_ADDRESS,
      hardhat: 0,
    },
  },
  networks: {
    hardhat: {
      gas: "auto",
      mining: {
        auto: true,
        interval: 5000,
        mempool: {
          order: "fifo",
        },
      },
      deploy: ["deploy/l1/"],
      saveDeployments: false,
    },
  },
};

const config = createConfig(baseConfig);
export default config;
