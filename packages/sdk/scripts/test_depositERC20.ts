import { ethers } from "ethers";
import { formatEther } from "ethers/lib/utils";
import { ServiceBridge } from '../src/service_bridge';

async function main() {


    // TODO: put me in a .env file
    const l1Url = 'http://localhost:8545'
    const l2Url = 'http://localhost:4011'


    const l1RpcProvider = new ethers.providers.JsonRpcProvider(l1Url)
    const l2RpcProvider = new ethers.providers.JsonRpcProvider(l2Url)
    const l1Wallet = new ethers.Wallet("0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6", l1RpcProvider)
    const l2Wallet = new ethers.Wallet('0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6', l2RpcProvider)

    const l1balance = await l1Wallet.getBalance()
    console.log("L1 balancee: ", l1balance)

    const serviceBridge = new ServiceBridge({
        l1SignerOrProvider: l1Wallet, // l1 signer
        l2SignerOrProvider: l2Wallet, // l2 signer
        l1ChainId: 1337,
        l2ChainId: 13527,
    });

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });