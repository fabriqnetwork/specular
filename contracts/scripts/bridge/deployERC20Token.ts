import { ethers } from "ethers";

type TokenInfo = {
  l1Address: String;
  name: String;
  symbol: String;
};

const address = "";
const abi = [
  "function createMintableERC20(address _remoteToken, string memory _name, string memory _symbol)",
  "event MintableERC20Created(address indexed localToken, address indexed remoteToken, address deployer)",
];

export async function deployERC20Token(
  token: TokenInfo,
  signer: ethers.Signer,
): Promise<String> {
  const mintableFactory = new ethers.Contract(address, abi, signer);
  const tx = await mintableFactory.createMintableERC20(
    token.l1Address,
    token.name,
    token.symbol,
  );

  const txWithLogs = await tx.wait();
  const deployEvent = mintableFactory.interface.parseLog(txWithLogs.logs[0]);

  return deployEvent.args.localToken;
}
