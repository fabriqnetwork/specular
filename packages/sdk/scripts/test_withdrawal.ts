import { ethers } from "ethers";
import { formatEther } from "ethers/lib/utils";
import { ServiceBridge } from '../src/service_bridge';
import { delay } from "../src/utils";

async function main() {

    // TODO: put me in a .env file
    const l1Url = 'http://localhost:8545'
    const l2Url = 'http://localhost:4011'


    const l1RpcProvider = new ethers.providers.JsonRpcProvider(l1Url)
    const l2RpcProvider = new ethers.providers.JsonRpcProvider(l2Url)
    const l1Wallet = new ethers.Wallet("0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6", l1RpcProvider)
    const l2Wallet = new ethers.Wallet('0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6', l2RpcProvider)

    const l2balance = await l2Wallet.getBalance()
    console.log("L2 balancee: ", l2balance)

    const serviceBridge = new ServiceBridge({
        l1SignerOrProvider: l1Wallet, // l1 signer
        l2SignerOrProvider: l2Wallet, // l2 signer
        l1ChainId: 1337,
        l2ChainId: 13527,
    });


    const withdrawalETHResponse = await serviceBridge.withdrawETH(200);


    // // 2 block confirmations
    const withdrawalETHReceipt = await withdrawalETHResponse.wait(2);

    console.log(withdrawalETHReceipt)

    const messageStatus = await serviceBridge.getDepositStatus(withdrawalETHReceipt)

    while (!(messageStatus == 1)) {
        await delay(500);
        console.log("...Waiting for the TX to be ready for finalization...")
    }

    const finalizeWithdrawalResponse = await serviceBridge.finalizeWithdrawal(withdrawalETHReceipt);

    const finalizeWithdrawalReceipt = finalizeWithdrawalResponse.wait()

    console.log({ finalizeWithdrawalReceipt });

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
