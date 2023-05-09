// run with: npx hardhat run <this file> --network local
import dotenv from "dotenv";
import { keccak256 } from "ethereumjs-util";
import { ethers, upgrades } from "hardhat";

dotenv.config({ path: __dirname + "/../.env" });

const main = async () => {
  const proxyAddress = "0xff00000000000000000000000000000000000001";

  const implSlot = toEip1967Hash("eip1967.proxy.implementation");

  const provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:4011"
  );
  const sequencer = new ethers.Wallet(
    `0x${process.env.SEQUENCER_PRIVATE_KEY}`,
    provider
  );

  const PlaceholderFactory = await ethers.getContractFactory("UUPSPlaceholder");

  const implementationAddress = await provider.getStorageAt(
    proxyAddress,
    implSlot
  );
  console.log({ implementationAddress });

  const proxy = await upgrades.forceImport(proxyAddress, PlaceholderFactory);

  const updateTx = await upgrades.upgradeProxy(
    proxy.address,
    PlaceholderFactory
  );
  console.log({ updateTx });
};

function toEip1967Hash(label: string): string {
  const hash = keccak256(Buffer.from(label));
  const bigNumber = BigInt("0x" + hash.toString("hex")) - 1n;
  return "0x" + bigNumber.toString(16);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
