import { BigNumber } from "ethers";
import { ethers } from "ethers";
import { formatEther } from "ethers/lib/utils";

import { CrossChainMessenger } from "../cross-chain-messenger";

async function main() {

    // Local chain IDs
    // 1337 L1 
    // 13527 L2 

    // TODO: put me in a .env file
    const l1Url = 'http://l1-geth:8545'
    const l2Url = 'http://sp-geth:4011'


    const l1RpcProvider = new ethers.providers.JsonRpcProvider(l1Url)
    const l2RpcProvider = new ethers.providers.JsonRpcProvider(l2Url)


    const crossChainMessenger = new CrossChainMessenger({
        l1ChainId: 1337,
        l2ChainId: 13527,
        l1SignerOrProvider: l1RpcProvider,
        l2SignerOrProvider: l2RpcProvider,
    });



    const depositETHResponse = await crossChainMessenger.depositETH(200);
    // 2 block confirmations
    const depositETHReceipt = await depositETHResponse.wait(2);

    const finalizeReceipt = crossChainMessenger.finalizeDeposit(depositETHReceipt);


}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
