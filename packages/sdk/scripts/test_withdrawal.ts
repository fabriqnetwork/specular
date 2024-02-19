import { ethers } from "ethers";
import { formatEther } from "ethers/lib/utils";
import { ServiceBridge } from '../src/service_bridge';
import { delay } from "../src/utils";
import {
    L1Portal__factory,
} from "../src/types/contracts";
import { getL1ContractsByNetworkId } from "../src/utils/constants";


async function main() {

    // TODO: put me in a .env file
    const l1Url = 'http://localhost:8545'
    const l2Url = 'http://localhost:4011'
    const l1RpcProvider = new ethers.providers.JsonRpcProvider(l1Url)
    const l2RpcProvider = new ethers.providers.JsonRpcProvider(l2Url)
    const l1Wallet = new ethers.Wallet("0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6", l1RpcProvider)
    const l2Wallet = new ethers.Wallet('0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6', l2RpcProvider)

    const l2balance = await l2Wallet.getBalance()
    console.log("L2 wallet balance: ", l2balance.toBigInt())

    const serviceBridge = new ServiceBridge({
        l1SignerOrProvider: l1Wallet, // l1 signer
        l2SignerOrProvider: l2Wallet, // l2 signer
        l1ChainId: 1337,
        l2ChainId: 13527,
    });

    const L1PortalAddress = getL1ContractsByNetworkId(serviceBridge.l1ChainId).L1Portal.toString()
    const l1Portal = L1Portal__factory.connect(L1PortalAddress, serviceBridge.l1SignerOrProvider)


    // funding the L1Portal
    const donateTx = await l1Portal.donateETH({
        value: ethers.utils.parseEther("1"),
    });
    await donateTx.wait();

    const value = ethers.utils.parseEther("0.000001");

    const withdrawalETHResponse = await serviceBridge.withdrawETH(value);

    const withdrawalETHReceipt = await withdrawalETHResponse.wait();

    console.log('withdrawalETHReceipt: ', withdrawalETHReceipt)
    let messageStatus = await serviceBridge.getWithdrawalStatus(withdrawalETHReceipt)

    // while the message status is pending, keep waiting
    while (messageStatus == 0) {
        console.log("...Waiting for the TX to be ready for finalization...")
        await delay(5000);
        messageStatus = await serviceBridge.getWithdrawalStatus(withdrawalETHReceipt);
    }

    const finalizeWithdrawalResponse = await serviceBridge.finalizeWithdrawal(withdrawalETHReceipt);

    const finalizeWithdrawalReceipt = await finalizeWithdrawalResponse.wait()

    console.log({ finalizeWithdrawalReceipt });

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
