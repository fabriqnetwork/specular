import { useState } from 'react';
import { ethers } from 'ethers';
import { getStorageKey } from '../utils';
import type { PendingWithdrawal, MessageProof } from "../types";
import {
  IL1Portal__factory,
  L1Oracle__factory,
} from "../typechain-types";

import {
  CHIADO_NETWORK_ID,
  CHIADO_RPC_URL,
  SPECULAR_NETWORK_ID,
  SPECULAR_RPC_URL,
  L1PORTAL_ADDRESS,
  L2PORTAL_ADDRESS,
  L1ORACLE_ADDRESS,
  DEPOSIT_BALANCE_THRESHOLD,
} from "../constants";

import { NETWORKS } from '../chains';

interface Data {
  status: string;
  error?: string;
  data?: string;
}

interface wallet {
    address: string;
    chainId: number;
    provider: any;
  }



const INITIAL_DATA: Data = { status: 'pending' };

async function generateWithdrawProof(
  withdrawal: PendingWithdrawal
): Promise<MessageProof> {
  console.log("generating proof");
  const l2Provider = new ethers.providers.StaticJsonRpcProvider(SPECULAR_RPC_URL);
  const rawProof = await l2Provider.send(
    "eth_getProof",
    [
      L2PORTAL_ADDRESS,
      [getStorageKey(withdrawal.withdrawalHash)],
      ethers.utils.hexlify(withdrawal.l2BlockNumber),
    ]
  );
  return {
    accountProof: rawProof.accountProof,
    storageProof: rawProof.storageProof[0].proof,
  };
}


type SwitchChainFunction = (arg: string) => void;

function useFinalizeWithdraw(switchChain: SwitchChainFunction) {
  const [data, setData] = useState<Data>(INITIAL_DATA);

  const finalizeWithdraw = async (wallet: wallet, amount: ethers.BigNumberish, pendingWithdraw:PendingWithdrawal): Promise<void> => {

    if (!wallet) {
      setData({ status: 'failed', error: "Wallet doesn't exist" });
      return;
    }

    setData({ status: 'loading' });

    const provider = await wallet.provider
    const signer = await provider.getSigner();

    const l1Portal = IL1Portal__factory.connect(
      L1PORTAL_ADDRESS,
      signer,
    );
    try {
      const proof = await generateWithdrawProof(pendingWithdraw);
      console.log(proof)
      console.log(pendingWithdraw)

      if(pendingWithdraw.assertionID) {
        const tx = await l1Portal.finalizeWithdrawalTransaction(
          pendingWithdraw.withdrawalTx,
          pendingWithdraw.assertionID,
          proof.accountProof,
          proof.storageProof
        );
        setData({ status: 'pending', data: tx.hash });
        await tx.wait();
        setData({ status: 'successful', data: tx.hash });
      }
      else {
        throw console.error("assertionID not found");

      }
    } catch (e) {
      console.error(e);
    }
    switchChain(CHIADO_NETWORK_ID.toString());
  };

  return { finalizeWithdraw, data };
}

export default useFinalizeWithdraw;
