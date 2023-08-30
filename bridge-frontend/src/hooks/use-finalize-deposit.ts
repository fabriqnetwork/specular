import { useState } from 'react';
import { ethers } from 'ethers';
import { getStorageKey, requestFundDeposit } from '../utils';
import type { PendingDeposit, MessageProof } from "../types";
import {
  L2Portal__factory,
  L1Oracle__factory
} from "../typechain-types";
import type { JsonRpcProvider } from "@ethersproject/providers";
import {
  CHIADO_RPC_URL,
  SPECULAR_NETWORK_ID,
  SPECULAR_RPC_URL,
  L1PORTAL_ADDRESS,
  L2PORTAL_ADDRESS,
  L1ORACLE_ADDRESS,
  DEPOSIT_BALANCE_THRESHOLD,
} from "../constants";
import {
  NETWORKS
} from "../chains";
interface Data {
  status: string;
  error?: string;
  data?: string;
}
interface PendingData {
  status: string;
  data: PendingDeposit;
}
interface wallet {
    address: string;
    chainId: number;
    provider: any;
  }



const INITIAL_DATA: Data = { status: 'waiting' };

async function generateDepositProof(
  deposit: PendingDeposit
): Promise<MessageProof> {
  if (deposit.proofL1BlockNumber === undefined) {
    throw new Error("proofL1BlockNumber is undefined");
  }
  let rawProof = undefined;
  while (rawProof === undefined) {
    console.log("try");
    try {
      const l1Provider = new ethers.providers.StaticJsonRpcProvider(CHIADO_RPC_URL);
      console.log("generateDepositProof proofL1BlockNumber is",deposit.proofL1BlockNumber)
      console.log("generateDepositProof depositHash is",deposit.depositHash)
      rawProof = await (l1Provider as JsonRpcProvider).send(
        "eth_getProof",
        [
          L1PORTAL_ADDRESS,
          [getStorageKey(deposit.depositHash)],
          deposit.proofL1BlockNumber,
        ]
      );
    } catch (e) {
      console.log("got error");
      console.error(e);
    }
    console.log("passed");
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  return {
    accountProof: rawProof.accountProof,
    storageProof: rawProof.storageProof[0].proof,
  };
}
type SwitchChainFunction = (arg: string) => void;

function useFinalizeDeposit(switchChain: SwitchChainFunction) {
  const [data, setData] = useState<Data>(INITIAL_DATA);

  const finalizeDeposit = async (wallet: wallet, amount: ethers.BigNumberish, pendingDeposit:PendingData, setPendingDeposit:any): Promise<void> => {

    if(pendingDeposit.status==='finalized'){
      setPendingDeposit({ status: 'finalized', data: pendingDeposit.data})
      return;
    }
    switchChain(SPECULAR_NETWORK_ID.toString())
    if (!wallet) {
      setData({ status: 'failed', error: "Wallet doesn't exist" });
      return;
    }
    const l2Provider = new ethers.providers.StaticJsonRpcProvider(SPECULAR_RPC_URL)
    const l2Balance  = await l2Provider.getBalance(wallet.address);

    setData({ status: 'loading' });
    console.log("Finalizing with l2 banance"+l2Balance);
    const targetBalance = ethers.utils.parseEther(ethers.utils.formatUnits(l2Balance, NETWORKS[SPECULAR_NETWORK_ID].nativeCurrency.decimals));
    if (DEPOSIT_BALANCE_THRESHOLD.gt(targetBalance)) {
      console.log("Sending Request");
      try {
        const txHash = await requestFundDeposit(pendingDeposit.data);
        console.log("Success Transaction :"+txHash);
        setData({ status: 'successful', data: txHash });
      } catch (e) {
        console.error(e);
      }
      return;
    }

    const l1Oracle = L1Oracle__factory.connect(
      L1ORACLE_ADDRESS,
      l2Provider,
    );
    try {
      console.log("Before l1Oracle");
      var latestBlockNumber = await l1Oracle.blockNumber();
      console.log("After l1Oracle");
      pendingDeposit.data.proofL1BlockNumber = latestBlockNumber.toNumber();
      console.log("pendingDeposit.data is "+pendingDeposit.data+" & proofL1BlockNumber is "+pendingDeposit.data.proofL1BlockNumber+" & Deposit hash is "+pendingDeposit.data.depositHash);
      const proof = await generateDepositProof(pendingDeposit.data);
      console.log("accountProof is "+proof.accountProof);
      console.log("storageProof is "+proof.storageProof);
      console.log("Chain Id is: "+wallet.chainId)
      const provider = await wallet.provider;
      const signer = await provider.getSigner();
      const l2Portal = L2Portal__factory.connect(
        L2PORTAL_ADDRESS,
        signer
      );
      console.log("L2 Portal Connected")

      console.log("version is", pendingDeposit.data.depositTx.version)
      console.log("nonce is", pendingDeposit.data.depositTx.nonce)
      console.log("sender is", pendingDeposit.data.depositTx.sender)
      console.log("target is", pendingDeposit.data.depositTx.target)
      console.log("value is", pendingDeposit.data.depositTx.value)
      console.log("gasLimit is", pendingDeposit.data.depositTx.gasLimit)
      console.log("data is", pendingDeposit.data.depositTx.data)

      console.log("depositTx is", pendingDeposit.data.depositTx)

      const tx = await l2Portal.finalizeDepositTransaction(
        pendingDeposit.data.depositTx,
        proof.accountProof,
        proof.storageProof
      );

      console.log("tx is "+tx+" And tx.hash is "+tx.hash);
      setData({ status: 'pending', data: tx.hash });
      await tx.wait();
      setData({ status: 'successful', data: tx.hash });
      console.log("successful tx is "+tx+" And tx.hash is "+tx.hash);

    } catch (errorCatched) {
      console.log("Error Cached at finalizeDepositData "+errorCatched)
      const err: any = errorCatched;
      let error = 'Transaction failed.';
      if (err.code === -32603) {
        error = 'Transaction was not sent because of the low gas price. Try to increase it.';
      }
      setData({ status: 'failed', error });
      console.log("failed tx with error "+err);
    }
  };
  console.log("finalizeDeposit.data is "+data.data+" and status"+data.status)
  return { finalizeDeposit, data };
}

export default useFinalizeDeposit;
