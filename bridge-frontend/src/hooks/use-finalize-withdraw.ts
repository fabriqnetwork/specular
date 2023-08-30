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
  const l2Provider = new ethers.providers.StaticJsonRpcProvider(SPECULAR_RPC_URL);
  const rawProof = await l2Provider.send(
    "eth_getProof",
    [
      L2PORTAL_ADDRESS,
      [getStorageKey(withdrawal.withdrawalHash)],
      withdrawal.l2BlockNumber,
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

    switchChain(CHIADO_NETWORK_ID.toString());
    const provider = await wallet.provider
    const signer = await (provider as any).getSigner();

    const l1Portal = IL1Portal__factory.connect(
      L2PORTAL_ADDRESS,
      signer,
    );
    try {
      const proof = await generateWithdrawProof(pendingWithdraw);
      console.log(proof)
      console.log(pendingWithdraw)

      if(pendingWithdraw.assertionID) {

        let gas = await l1Portal.estimateGas.finalizeWithdrawalTransaction(
          pendingWithdraw.withdrawalTx,
          pendingWithdraw.assertionID,
          proof.accountProof,
          proof.storageProof,
          { gasLimit: 1000000 }, // avoid gas estimation error
        );
        console.log("gas", gas);
        gas = gas.add(150000); // extra gas to pass gas limit check in finalization
        const tx = await l1Portal.finalizeWithdrawalTransaction(
          pendingWithdraw.withdrawalTx,
          pendingWithdraw.assertionID,
          proof.accountProof,
          proof.storageProof,
          { gasLimit: gas },
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
