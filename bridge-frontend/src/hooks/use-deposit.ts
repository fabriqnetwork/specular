import { useState } from 'react';
import { ethers, BigNumberish } from 'ethers';

import {
  L1_BRIDGE_ADDR,
  L1PORTAL_ADDRESS
} from "../constants";
import {
  NETWORKS
} from "../chains";
import { TOKEN, erc20Abi } from '../tokens';

import {L1StandardBridge__factory } from "../typechain-types"

interface Data {
  status: string;
  error?: string;
  data?: ethers.providers.TransactionResponse;
}

interface wallet {
    address: string;
    chainId: number;
    provider: any;
  }



const INITIAL_DATA: Data = { status: 'waiting' };

function useDeposit() {
  const [data, setData] = useState<Data>(INITIAL_DATA);

  const deposit = async (wallet: wallet, amount: ethers.BigNumberish, selectedTokenKey: number): Promise<void> => {
    console.log("From wallet " + wallet.address + " Amount " + amount);

    if (!wallet) {
      setData({ status: 'failed', error: "Wallet doesn't exist" });
      return;
    }

    try {
      const selectedToken = TOKEN[selectedTokenKey];
      const signer = await (wallet.provider as any).getSigner();
      const l1StandardBridge = L1StandardBridge__factory.connect(L1_BRIDGE_ADDR ,signer)

      let tx;
      if(selectedToken.l1TokenContract===""){

        tx = await l1StandardBridge.bridgeETH(200_000, [], {
          value: amount,
        });
      } else{
        console.log("erc20");
        const l1Token = new ethers.Contract(
          selectedToken.l1TokenContract,
          erc20Abi,
          signer
        );
        const approveTx = await l1Token.approve(
          L1_BRIDGE_ADDR,
          amount
        );
        await approveTx.wait();

        tx = await l1StandardBridge.bridgeERC20(
          selectedToken.l1TokenContract,
          selectedToken.l2TokenContract,
          amount,
          200_000,
          []
        );
    }

      console.log(tx)
      setData({ status: 'pending', data: tx });
      await tx.wait();
      setData({ status: 'successful', data: tx });

    } catch (errorCatched) {
      const err: any = errorCatched;
      let error = 'Transaction failed.';
      if (err.code === -32603) {
        error = 'Transaction was not sent because of the low gas price. Try to increase it.';
      }
      setData({ status: 'failed', error });
      console.log(err);
    }
  };

  const resetData = () => {
    setData(INITIAL_DATA);
  };

  return { deposit, data, resetData };
}

export default useDeposit;
