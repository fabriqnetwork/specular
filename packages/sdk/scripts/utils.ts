import { ethers } from "ethers";
import { formatEther } from "ethers/lib/utils";
import { ServiceBridge } from "../src/service_bridge";
import { delay } from "../src/utils";
import { resolve } from "path";
import dotenv from "dotenv";
dotenv.config({ path: resolve(__dirname, "../.env.test") });

async function main() {
  const l1Url: string = process.env.L1_URL!;
  const l2Url: string = process.env.L2_URL!;
  const l1ChainId: number = parseInt(process.env.L1_CHAIN_ID!);
  const l2ChainId: number = parseInt(process.env.L2_CHAIN_ID!);

  const l1RpcProvider = new ethers.providers.JsonRpcProvider(l1Url);
  const l2RpcProvider = new ethers.providers.JsonRpcProvider(l2Url);
  const l1Wallet = new ethers.Wallet(
    "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
    l1RpcProvider,
  );
  const l2Wallet = new ethers.Wallet(
    "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
    l2RpcProvider,
  );

  const l1balance = await l1Wallet.getBalance();
  console.log("L1 balancee: ", l1balance);

  const serviceBridge = new ServiceBridge({
    l1SignerOrProvider: l1Wallet, // l1 signer
    l2SignerOrProvider: l2Wallet, // l2 signer
    l1ChainId,
    l2ChainId,
  });

  const test = await serviceBridge.getLastTransactionsFromAddress(
    l1Wallet.address,
    10,
  );

  console.log(l1Wallet.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
