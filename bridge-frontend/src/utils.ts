import { ethers } from "ethers";
import type { PendingDeposit, PendingWithdrawal } from "./types";
import { RELAYER_ENDPOINT } from "./constants";

export function getStorageKey(messageHash: string) {
  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      ["bytes32", "uint256"],
      [messageHash, 0]
    )
  );
}

export async function requestFundDeposit(deposit: PendingDeposit): Promise<string> {
  const reqBody = {
    nonce: deposit.depositTx.nonce.toString(),
    sender: deposit.depositTx.sender,
    target: deposit.depositTx.target,
    value: deposit.depositTx.value.toString(),
    gasLimit: deposit.depositTx.gasLimit.toString(),
    data: deposit.depositTx.data,
    depositHash: deposit.depositHash,
  };
  const reqOpt = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(reqBody),
  };
  const res = await fetch(`${RELAYER_ENDPOINT}/fundDeposit`, reqOpt);
  if (!res.ok) {
    console.error(res);
    throw new Error(`Failed to request fund deposit: ${res.statusText}`);
  }
  const resBody = await res.json();
  return resBody["txHash"];
}

export async function requestFundWithdraw(deposit: PendingWithdrawal): Promise<string> {
  const reqBody = {
    nonce: deposit.withdrawalTx.nonce.toString(),
    sender: deposit.withdrawalTx.sender,
    target: deposit.withdrawalTx.target,
    value: deposit.withdrawalTx.value.toString(),
    gasLimit: deposit.withdrawalTx.gasLimit.toString(),
    data: deposit.withdrawalTx.data,
    depositHash: deposit.withdrawalHash,
  };
  const reqOpt = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(reqBody),
  };
  const res = await fetch(`${RELAYER_ENDPOINT}/fundDeposit`, reqOpt);
  if (!res.ok) {
    console.error(res);
    throw new Error(`Failed to request fund deposit: ${res.statusText}`);
  }
  const resBody = await res.json();
  return resBody["txHash"];
}
