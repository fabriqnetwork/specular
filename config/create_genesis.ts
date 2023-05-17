import { BigNumber, ethers } from "ethers";
import { exec } from "child_process";
import { glob } from "glob";
import { keccak256 } from "ethereumjs-util";
import hre from "hardhat";
import fs from "fs";
import path from "path";
import util from "node:util";

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
  let artifacts;
  try {
    artifacts = await hre.artifacts.readArtifact(contractName);
  } catch (error) {
    console.warn(
      `could not find artifacts for ${contractName}, checking in node_modules/@openzeppelin`
    );
    const paths =
      await glob(`../contracts/node_modules/@openzeppelin/**/build/**/${contractName}.json`);

    if (paths.length !== 1) throw Error("no unique artifacts found");

    artifacts = JSON.parse(fs.readFileSync(paths[0], "utf-8"));
  }

  return artifacts;
}

async function getStorageLayout(contractName: string, artifact: any) {
  // unfortunately this is not exposed through the hre so we have to get it manually
  // this is following the approach taken by OZ in the HH upgrade plugin, see:
  // https://github.com/OpenZeppelin/openzeppelin-upgrades/blob/1e28bce2b6bae17c024350c98c6e0c511f3091d3/packages/core/src/validate/query.ts#L81
  const validationsPath = path.join(hre.config.paths.cache, "validations.json");
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
  let storageLayout;

  if (p.balance) data.set("balance", p.balance);

  if (p.contract) {
    artifact = await getArtifact(p.contract);
    storageLayout = await getStorageLayout(p.contract, artifact);
    data.set("code", artifact.deployedBytecode);
  }

  if (!storageLayout && p.contract) {
    console.warn(
      `could not get storage layout for ${p.contract}, please ensure this is intentional`
    );
    ;
  }

  const storage = new Map();
  if (storageLayout) {
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
  }

  if (p.storage) {
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
