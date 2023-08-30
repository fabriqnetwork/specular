import { useState } from 'react';
import { ethers,  BytesLike } from 'ethers';
import { getStorageKey, requestFundDeposit } from '../utils';
import type { PendingDeposit, MessageProof } from "../types";
import {IL2Portal__factory, L1Oracle__factory} from "../typechain-types"

import type { JsonRpcProvider } from "@ethersproject/providers";
import {
  CrossDomainMessage,
  Data
} from "../types";
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

interface wallet {
    address: string;
    chainId: number;
    provider: any;
  }




type SwitchChainFunction = (arg: string) => void;

function useFinalizeDeposit(switchChain: SwitchChainFunction) {

  const finalizeDeposit = async (wallet: wallet, amount: ethers.BigNumberish, depositData:Data, setPendingDeposit:any): Promise<void> => {
    console.log("Finalize Chain Id is: "+wallet.chainId)

    const signer = await (wallet.provider as any).getSigner();
    console.log("Before l1Oracle")
    const l1Oracle = L1Oracle__factory.connect("0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1", signer)
    console.log("Before l2Portal")
    const l2Portal = IL2Portal__factory.connect(
        "0x8438Ad1C834623CfF278AB6829a248E37C2D7E3f",
        signer
      );
    const l1Provider = new ethers.providers.StaticJsonRpcProvider(CHIADO_RPC_URL);
    const latestBlockNumber = await l1Oracle.blockNumber()

    const rawBlock = await l1Provider.send("eth_getBlockByNumber", [
      ethers.utils.hexValue(latestBlockNumber),
      false, // We only want the block header
    ]);
    console.log("rawBlock",rawBlock)
    const stateRoot = l1Provider.formatter.hash(rawBlock.stateRoot);
    console.log("stateRoot",stateRoot)

    // const proof = await l1Provider.send("eth_getProof", [
    //   l1Portal.address,
    //   [getStorageKey(tx.hash)],
    //   "latest",
    // ]);
    // console.log("proof",proof)
    // const accountProof = proof.accountProof;
    // console.log("accountProof",accountProof)
    // const storageProof = proof.storageProof[0].proof;
    // console.log("storageProof",storageProof)



    let finalizeTx;
    // console.log("l1Oracle set")
    // const gasPrice = ethers.utils.parseUnits("2499", "gwei"); // 20 Gwei, for example

    // await l1Oracle.setL1OracleValues(depositData.l1BlockNumber || "", depositData.stateRoot || "", 0,
    // {
    //   gasPrice: gasPrice
    // });
    // console.log("l1Oracle",depositData.l1BlockNumber || "",depositData.stateRoot || "")
    // if(depositData.crossDomainMessage){
    //   console.log("crossDomainMessage",depositData.crossDomainMessage)
    //   console.log("accountProof",depositData.accountProof)
    //   console.log("storageProof",depositData.crossDomainMessage)
    // const finalizeTx = await l2Portal.finalizeDepositTransaction(
    //   depositData.crossDomainMessage,
    //   depositData.accountProof,
    //   depositData.storageProof
    // );
  //   console.log("finalizeTx",finalizeTx)
  //   await finalizeTx.wait();
  //   console.log("finalizeTx2",finalizeTx)
  // } else{
    console.log("Invalid Cross Domain Message")
  // }
  console.log("finalizeTx2",finalizeTx)
  };
  return { finalizeDeposit};
}

export default useFinalizeDeposit;
