require('dotenv').config()
require('@nomiclabs/hardhat-ethers')
require('@nomiclabs/hardhat-etherscan')
require('@nomiclabs/hardhat-waffle')
require('hardhat-abi-exporter')
require('hardhat-deploy')
require('@openzeppelin/hardhat-upgrades')

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: {
    version: '0.8.4',
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  },
  paths: {
    sources: './src'
  },
  abiExporter: {
    path: './abi',
    runOnCompile: true,
    clear: true
  },
  namedAccounts: {
    sequencer: 0,
    validator: 1
  },
  networks: {
    hardhat: {
      gas: 'auto',
      mining: {
        auto: true,
        interval: 5000,
        mempool: {
          order: 'fifo'
        }
      }
    }
  }
}
