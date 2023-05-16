import { BigNumber, ethers } from "ethers";
import { exec } from "child_process";
import { glob } from "glob";
import { keccak256 } from "ethereumjs-util";
import fs from "fs";
import path from "path";
import util from "node:util";

const CONTRACTS_PATH = "../contracts";

type PreDeploy = {
  address: string;
  balance: string;
  contract: string | undefined;
  storage: { [key: string]: string };
};

async function main() {
  const { baseGenesisPath, genesisPath } = parseArgs();
  await generateGenesisFile(baseGenesisPath, genesisPath);
}

export async function generateGenesisFile(
  baseGenesisPath: string,
  genesisPath: string
) {
  const baseGenesis = JSON.parse(fs.readFileSync(baseGenesisPath, "utf-8"));

  const alloc = new Map();
  await Promise.all(
    baseGenesis.preDeploy.map((p: any) => parsePreDeploy(p, alloc))
  );

  baseGenesis.preDeploy = undefined;
  baseGenesis.alloc = Object.fromEntries(alloc);

  fs.writeFileSync(genesisPath, JSON.stringify(baseGenesis, null, 2), "utf-8");
}

function parseArgs() {
  const inFlagIndex = process.argv.indexOf("--in");
  let baseGenesisPath;

  if (inFlagIndex > -1) {
    baseGenesisPath = process.argv[inFlagIndex + 1];
  } else {
    throw Error("Please specify the base genesis path");
  }

  const outFlagIndex = process.argv.indexOf("--out");
  let genesisPath;

  if (outFlagIndex > -1) {
    genesisPath = process.argv[outFlagIndex + 1];
  } else {
    console.log("Setting out genesis path same as base genesis path");
    genesisPath = path.join(path.dirname(baseGenesisPath), "genesis.json");
  }

  return { baseGenesisPath, genesisPath };
}

async function getArtifact(contractName: string) {
  const artifactPath = await glob(`${CONTRACTS_PATH}/artifacts/**/${contractName}.json`);

  if (artifactPath.length !== 1) {
    throw Error(`could not find unique artifact for ${contractName}`);
  }

  return JSON.parse(fs.readFileSync(artifactPath[0], "utf-8"));
}

async function getStorageLayout(contractName: string, artifact: any) {
  const validationsPath = `${CONTRACTS_PATH}/cache/validations.json`;
  const validations = JSON.parse(fs.readFileSync(validationsPath, "utf-8"));

  // the version of the validation is the keccak hash of the contracts bytecode
  const buf = Buffer.from(artifact.bytecode.replace(/^0x/, ''), 'hex');
  const version = keccak256(buf).toString('hex');

  let storageLayout;
  let fullContractName;

  for (const validation of validations.log) {
    fullContractName = Object.keys(validation).find(
      name => validation[name].version?.withMetadata === version,
    );
    if (fullContractName !== undefined) {
      storageLayout = validation[fullContractName].layout.storage;
      break;
    }
  }

  return storageLayout;
}

async function parsePreDeploy(p: PreDeploy, alloc: any) {
  const execPromise = util.promisify(exec);
  const data = new Map();
  let artifact;

  if (p.balance) data.set("balance", p.balance);

  if (p.contract) {
    artifact = await getArtifact(p.contract);
    data.set("code", artifact.deployedBytecode);
  }

  if (p.storage && p.contract) {
    const storageLayout = await getStorageLayout(p.contract, artifact);

    const storage = new Map();
    for (const s of storageLayout) {
      if (!p.storage[s.label]) continue;

      const slot = ethers.utils.hexZeroPad(
        BigNumber.from(s.slot).toHexString(),
        32
      );

      let value = new Uint8Array(32);
      if (storage.has(slot)) {
        value = storage.get(slot);
      }

      const newStorage = ethers.utils.arrayify(
        BigNumber.from(p.storage[s.label]).toHexString()
      );

      for (let i = 0; i < newStorage.length; i++) {
        value[value.length - s.offset - newStorage.length + i] = newStorage[i];
      }

      storage.set(slot, ethers.utils.hexZeroPad(value, 32));
    }

    for (const s of Object.entries(p.storage)) {
      if (!ethers.utils.isHexString(s[0])) continue;

      const slot = ethers.utils.hexZeroPad(s[0], 32);
      const value = ethers.utils.hexZeroPad(s[1], 32);
      storage.set(slot, value);
    }

    data.set("storage", Object.fromEntries(storage));
  }

  alloc.set(p.address, Object.fromEntries(data));
}

if (!require.main!.loaded) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
