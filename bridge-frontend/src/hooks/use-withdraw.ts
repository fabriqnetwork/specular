import { useState } from 'react';
import { ethers, BigNumberish } from 'ethers';

import {
  L2PORTAL_ADDRESS,
} from "../constants";
import {
  NETWORKS
} from "../chains";

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

function weiToEther(wei: BigNumberish): string {
  const weiPerEther: ethers.BigNumber = ethers.BigNumber.from("1000000000000000000"); // 1 ether = 10^18 wei
  const weiValue: ethers.BigNumber = ethers.BigNumber.from(wei);

  const etherValue: ethers.BigNumber = weiValue.div(weiPerEther);
  const remainder: ethers.BigNumber = weiValue.mod(weiPerEther);

  const formattedEther: string = etherValue.toString();
  const formattedRemainder: string = remainder.toString().padStart(18, "0"); // Pad with leading zeros if necessary

  return `${formattedEther}.${formattedRemainder}`;
}

function useWithdraw() {
  const [data, setData] = useState<Data>(INITIAL_DATA);

  const withdraw = async (wallet: wallet, amount: ethers.BigNumberish): Promise<void> => {
    console.log("From wallet " + wallet.address + " Amount " + amount);
    setData({ status: 'loading' });

    if (!wallet) {
      setData({ status: 'failed', error: "Wallet doesn't exist" });
      return;
    }

    try {
      console.log("In Try Block with amount " + weiToEther(amount));

      const signer = await (wallet.provider as any).getSigner();
      const tx = await signer.sendTransaction({
        to: L2PORTAL_ADDRESS,
        value: ethers.utils.parseUnits(weiToEther(amount),  NETWORKS[wallet.chainId].nativeCurrency.decimals),
      });

      console.log(tx)

      await tx.wait();
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
