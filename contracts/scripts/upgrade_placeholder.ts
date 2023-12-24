// run with: npx hardhat run <this file> --network local
import dotenv from "dotenv";
import { keccak256 } from "ethereumjs-util";
import { ethers, upgrades } from "hardhat";

dotenv.config({ path: __dirname + "/../.env" });

const main = async () => {
  const UUPSPlaceholderFactory = await ethers.getContractFactory(
    "UUPSPlaceholder"
  );
  const FaucetFactory = await ethers.getContractFactory("Faucet");

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

  const upgraded = await upgrades.upgradeProxy(proxy.address, FaucetFactory);
  console.log({ upgraded });

  const tx = await upgraded.owner();
  console.log({ tx });
};

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
