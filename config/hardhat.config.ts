import { HardhatUserConfig } from "hardhat/types";

const config: HardhatUserConfig = {
  solidity: "0.8.17",
  paths: {
    artifacts: "../contracts/artifacts/",
    cache: "../contracts/cache/"
  }
};

export default config;
