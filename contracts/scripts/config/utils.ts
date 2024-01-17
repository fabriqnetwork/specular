import { ethers } from "ethers";

export function parseFlag(flag: string, defaultValue?: string): string {
  const flagIndex = process.argv.indexOf(flag);
  let value: string | undefined;

  if (flagIndex < 0) {
    value = defaultValue;
  } else {
    value = process.argv[flagIndex + 1];
  }

  if (value === undefined) throw Error(`no value set for "${flag}"`);
  return value;
}

export function numberStrToPaddedHex(numStr: string, length: number): string {
  return ethers.utils.hexZeroPad(
    ethers.BigNumber.from(numStr).toHexString(),
    length,
  );
}
