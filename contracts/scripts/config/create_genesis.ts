import { BigNumber, ethers } from "ethers";
import { glob } from "glob";
import { keccak256 } from "ethereumjs-util";
import hre from "hardhat";
import fs from "fs";
import path from "path";
import { parseFlag } from "./utils";

type PreDeploy = {
  address: string;
  balance: string;
  contract: string | undefined;
  storage: { [key: string]: string };
};

async function main() {
  const baseGenesisPath = parseFlag("--in")
  const defaultGenesisPath = path.join(path.dirname(baseGenesisPath), "genesis.json");
  const genesisPath = parseFlag("--out", defaultGenesisPath)
  const l1RpcUrl = parseFlag("--l1-rpc-url")
  await generateGenesisFile(baseGenesisPath, genesisPath, l1RpcUrl);
}

export async function generateGenesisFile(
  baseGenesisPath: string,
  genesisPath: string,
  l1RpcUrl: string
) {
  const baseGenesis = JSON.parse(fs.readFileSync(baseGenesisPath, "utf-8"));

  const alloc = new Map();
  await Promise.all(
    baseGenesis.preDeploy.map((p: any) => parsePreDeploy(p, alloc))
  );

  baseGenesis.preDeploy = undefined;
  baseGenesis.alloc = Object.fromEntries(alloc);

  // ethers v6 provides a better handle to close the websocket,
  // but we need to do this in v5 so the script terminates
  try {
    const provider = new ethers.providers.WebSocketProvider(l1RpcUrl)
    const block = await provider.getBlock("safe")
    baseGenesis.timestamp = block.timestamp
    provider._websocket.terminate();
  } catch (error) {
    console.error(`could not get l1 safe block from network: ${l1RpcUrl}, error: ${error}`);
  }

  fs.writeFileSync(genesisPath, JSON.stringify(baseGenesis, null, 2), "utf-8");
  console.log(`successfully wrote genesis file to: ${genesisPath}`)
}

/**
 * Reads the compilation artifact of a contract compiled with hardhat
 * If no artifact is found, the function will look in the artifacts shipped with the @openzeppelin package
 * @param {string} contractName - unique name of the contract
 */
async function getArtifact(contractName: string) {
  let artifacts;
  try {
    artifacts = await hre.artifacts.readArtifact(contractName);
  } catch (error) {
    console.warn(
      `could not find artifacts for ${contractName}, checking in node_modules/@openzeppelin`
    );
    const paths =
      await glob(`./node_modules/@openzeppelin/**/build/**/${contractName}.json`);

    if (paths.length !== 1) throw Error("no unique artifacts found");

    artifacts = JSON.parse(fs.readFileSync(paths[0], "utf-8"));
  }

  return artifacts;
}

/**
 * Reads the storage layout from a contract compiled with hardhat
 * unfortunately this is not exposed through the hre directly so we have to get it manually
 * this is following the approach taken by OZ in the HH upgrade plugin, see:
 * https://github.com/OpenZeppelin/openzeppelin-upgrades/blob/1e28bce2b6bae17c024350c98c6e0c511f3091d3/packages/core/src/validate/query.ts#L81
 * @param {any} artifact - the contracts compilation artifact
 */
async function getStorageLayout(artifact: any) {
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

/**
 * Turn a PreDeploy object into a well formated genesis.json alloc entry
 * @param {PreDeploy} p - the configuration of the PreDeploy
 * @param {any} alloc - a map containing all alloc entries to be written in the genesis file
 */
async function parsePreDeploy(p: PreDeploy, alloc: any) {
  if (alloc.has(p.address)) {
    throw Error(`multiple pre-deploys specified for address: ${p.address}`);
  }

  const data = new Map();
  let artifact;
  let storageLayout;

  if (p.balance) data.set("balance", p.balance);

  if (p.contract) {
    artifact = await getArtifact(p.contract);
    storageLayout = await getStorageLayout(artifact);
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
      if (!p.storage || !p.storage[s.label]) continue;

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
