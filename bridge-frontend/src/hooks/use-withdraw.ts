import { useState } from 'react';
import { ethers, BigNumberish } from 'ethers';

import {
  L2PORTAL_ADDRESS,
  L2_BRIDGE_ADDR
} from "../constants";
import {
  NETWORKS
} from "../chains";
import { TOKEN, erc20Abi } from '../tokens';

import {L2StandardBridge__factory } from "../typechain-types"

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



const INITIAL_DATA: Data = { status: 'pending' };

function useWithdraw() {
  const [data, setData] = useState<Data>(INITIAL_DATA);

  const withdraw = async (wallet: wallet, amount: ethers.BigNumberish, selectedTokenKey: number): Promise<void> => {
    console.log("From wallet " + wallet.address + " Amount " + amount);
    setData({ status: 'loading' });

    if (!wallet) {
      setData({ status: 'failed', error: "Wallet doesn't exist" });
      return;
    }

    try {
      const selectedToken = TOKEN[selectedTokenKey];
      const signer = await (wallet.provider as any).getSigner();
      const l2StandardBridge = L2StandardBridge__factory.connect(L2_BRIDGE_ADDR ,signer)

      let tx;
      if(selectedToken.l1TokenContract===""){
        // tx = await l2StandardBridge.bridgeETH(200_000, [], {
        //   value: amount,
        // });

        tx = await signer.sendTransaction({
          to: L2PORTAL_ADDRESS,
          value: amount,
        });
      } else{
        console.log("erc20");
        const l2Token = new ethers.Contract(
          selectedToken.l2TokenContract,
          erc20Abi,
          signer
        );
        const approveTx = await l2Token.approve(
          L2_BRIDGE_ADDR,
          amount
        );
        await approveTx.wait();

        tx = await l2StandardBridge.bridgeERC20(
          selectedToken.l2TokenContract,
          selectedToken.l1TokenContract,
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

  return { withdraw, data, resetData };
}

export default useWithdraw;
