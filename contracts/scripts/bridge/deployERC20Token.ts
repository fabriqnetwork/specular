import { ethers } from "ethers";
import { env } from "process";
import * as dotenv from "dotenv";
import { program } from "commander";

type TokenInfo = {
  l1Address: String;
  name: String;
  symbol: String;
};

const address = "0x2A000000000000000000000000000000000000F0";
const abi = [
  "function createMintableERC20(address _remoteToken, string memory _name, string memory _symbol)",
  "event MintableERC20Created(address indexed localToken, address indexed remoteToken, address deployer)",
];

export async function deployERC20Token(
  token: TokenInfo,
  signer: ethers.Signer,
): Promise<String> {
  const mintableFactory = new ethers.Contract(address, abi, signer);

  console.log({ mintableFactory });
  console.log({ token });

  const tx = await mintableFactory.createMintableERC20(
    token.l1Address,
    token.name,
    token.symbol,
  );
  console.log({ tx });

  const txWithLogs = await tx.wait();
  const deployEvent = mintableFactory.interface.parseLog(txWithLogs.logs[0]);

  return deployEvent.args.localToken;
}

program
  .requiredOption("--rpc <string>")
  .requiredOption("--name <string>")
  .requiredOption("--symbol <string>")
  .requiredOption("--address <string>");

async function main() {
  dotenv.config();
  program.parse(process.argv);
  const options = program.opts();
  if (!env.DEPLOYER_PRIVATE_KEY) throw "please provide deployer private key";

  const key = new ethers.utils.SigningKey(env.DEPLOYER_PRIVATE_KEY);
  const provider = new ethers.providers.JsonRpcProvider(options.rpc);
  const signer = new ethers.Wallet(key, provider);
  const token = {
    l1Address: options.address,
    name: options.name,
    symbol: options.symbol,
  };

  const l2address = await deployERC20Token(token, signer);
  console.log("token deployed on L2 at ", l2address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
