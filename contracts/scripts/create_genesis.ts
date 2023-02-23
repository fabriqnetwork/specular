import { BigNumber, ethers } from "ethers";
import fs from "fs";
import FaucetJson from "../artifacts/src/pre-deploy/Faucet.sol/Faucet.json";
import assert from "assert";
let GenesisJson;

interface contractObject {
  code: string;
  balance: string;
  storage: any;
}

const createContractObject = (
  deployedBytecode: string,
  contractBalance: BigNumber,
  storageSlots: Array<string>,
  valueAtSlots: Array<string>
): contractObject => {
  assert(
    storageSlots.length == valueAtSlots.length,
    "incorrect storage-values array lengths"
  );

  const storageSlotsObj: any = {};
  for (let i = 0; i < storageSlots.length; i++) {
    storageSlotsObj[storageSlots[i].toString()] = valueAtSlots[i];
  }

  const contractObject = {
    code: deployedBytecode,
    balance: Number(contractBalance).toString(),
    storage: storageSlotsObj,
  };
  return contractObject;
};

const createFaucetContractObject = (): contractObject => {
  const faucetDeployedBytecode = FaucetJson.deployedBytecode;
  const faucetBalance = ethers.BigNumber.from("10").pow(20);

  const storageSlots = [];
  const valueAtSlots = [];

  // Storage Slot 0 stores the address of the owner
  storageSlots[0] =
    "0x0000000000000000000000000000000000000000000000000000000000000000";
  // Storage Slot 1 stores the amountAllowed
  storageSlots[1] =
    "0x0000000000000000000000000000000000000000000000000000000000000001";

  // Owner Address - 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266
  valueAtSlots[0] =
    "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266";
  // Amount Allowed - 0.01 ETH
  valueAtSlots[1] =
    "0x000000000000000000000000000000000000000000000000002386f26fc10000";

  const faucetObject = createContractObject(
    faucetDeployedBytecode,
    faucetBalance,
    storageSlots,
    valueAtSlots
  );
  return faucetObject;
};

const main = () => {
  const inFlagIndex = process.argv.indexOf("--in");
  let baseGenesisPath;

  if (inFlagIndex > -1) {
    baseGenesisPath = process.argv[inFlagIndex + 1];
    GenesisJson = require(baseGenesisPath);
  } else {
    throw new Error("Please specify the base genesis path");
  }

  const outFlagIndex = process.argv.indexOf("--out");
  let genesisPath;

  if (outFlagIndex > -1) {
    genesisPath = process.argv[outFlagIndex + 1];
  } else {
    console.log("Setting out genesis path same as base genesis path");
    genesisPath = baseGenesisPath;
  }

  // Address the faucet will be deployed to - address(20)
  const faucetAddress = "0x0000000000000000000000000000000000000020";

  GenesisJson.alloc[faucetAddress.toString()] = createFaucetContractObject();
  fs.writeFileSync(genesisPath, JSON.stringify(GenesisJson, null, 2));
};

main();
