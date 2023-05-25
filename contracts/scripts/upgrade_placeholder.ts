// run with: npx hardhat run <this file> --network local
import dotenv from "dotenv";
import { keccak256 } from "ethereumjs-util";
import { ethers, upgrades } from "hardhat";

dotenv.config({ path: __dirname + "/../.env" });

const main = async () => {
  const UUPSPlaceholderFactory = await ethers.getContractFactory(
    "UUPSPlaceholder"
  );
  const V2Factory = await ethers.getContractFactory("UUPSPlaceholderV2");

  const proxyAddress = "0xff00000000000000000000000000000000000001";

  console.log(
    "Implementation address: " +
      (await upgrades.erc1967.getImplementationAddress(proxyAddress))
  );
  console.log(
    "Admin address: " + (await upgrades.erc1967.getAdminAddress(proxyAddress))
  );

  const proxy = await upgrades.forceImport(
    proxyAddress,
    UUPSPlaceholderFactory
  );

  const init = await proxy.initialize();
  console.log({ init });

  const upgraded = await upgrades.upgradeProxy(proxy.address, V2Factory);
  console.log({ upgraded });

  // this function is only availabe in the V2 contract, so this serves as a simple test
  // to check if the upgrade was successful
  const greeting = await upgraded.greet();
  console.log({ greeting });
};

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
