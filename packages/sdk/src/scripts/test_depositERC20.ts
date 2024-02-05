import { ethers } from "ethers";
import { formatEther } from "ethers/lib/utils";

import { CrossChainMessenger } from "../cross-chain-messenger";

async function main() {


    // TODO: put me in a .env file
    const l1Url = 'http://l1-geth:8545'
    const l2Url = 'http://sp-geth:4011'


    const l1RpcProvider = new ethers.providers.JsonRpcProvider(l1Url)
    const l2RpcProvider = new ethers.providers.JsonRpcProvider(l2Url)

    console.log(l1RpcProvider)

    const crossChainMessenger = new CrossChainMessenger({
        l1ChainId: 1337,
        l2ChainId: 13527,
        l1SignerOrProvider: l1RpcProvider,
        l2SignerOrProvider: l2RpcProvider,
    });

    // if local env 
    // deploy token on l1
    // deploy token on l2

    // approve L1 token
    // approve L2 token

    // const depositETHResponse = await crossChainMessenger.depositERC20(200);
    // 2 block confirmations
    // const depositETHReceipt = await depositETHResponse.wait(2);

    // const finalizeReceipt = crossChainMessenger.finalizeDeposit(depositETHReceipt);


}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
