import "dotenv/config";
import "@nomiclabs/hardhat-ethers";
import {HardhatUserConfig} from 'hardhat/types';
import "@nomiclabs/hardhat-etherscan";
import "@nomiclabs/hardhat-waffle";
import 'hardhat-abi-exporter';
import 'hardhat-deploy';
import '@openzeppelin/hardhat-upgrades';

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
const config: HardhatUserConfig = {
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
      "goerli": '0xf112347faDA222A95d84626b19b2af1DB6672C18'
    },
    validator: 1,
  },
  networks: {
    goerli: {
      url: `https://goerli.infura.io/v3/${process.env.INFURA_KEY}`,
      chainId: 5,
      accounts: [`0x${process.env.SEQUENCER_PRIVATE_KEY}`],
    },
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