import { ethers } from "ethers";

export function getStorageKey(messageHash: string) {
  return ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      ["bytes32", "uint256"],
      [messageHash, 0]
    )
  );
}
