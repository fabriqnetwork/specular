import { BigNumber, ethers } from "ethers";
import { exec } from "child_process";
import { glob } from "glob";
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
  const artifactPath = await glob(`**/artifacts/**/${contractName}.json`, {
    ignore: "node_modules/**",
  });

  if (artifactPath.length !== 1) {
    throw Error(`could not find unique artifact for ${contractName}`);
  }

  return JSON.parse(fs.readFileSync(artifactPath[0], "utf-8"));
}

// TODO: this is not very efficient,
// we should be able to narrow down the search once we know
// where the file will live / how it will be run
async function getStorageLayout(contractName: string) {
  // we can read the storage layout from the HH artifacts like this:
  // 1. read the contracts debug file to get the build-info hash of the most recent compilation
  // 2. get the full source path of the contract
  // 3. use the full source path to look up the storage layout in the build-info file

  const debugFile = await glob(`**/artifacts/**/${contractName}.dbg.json`, {
    ignore: "node_modules/**",
  });

  if (debugFile.length !== 1) {
    throw Error(`could not find unique build info for ${contractName}`);
  }

  const debug = JSON.parse(fs.readFileSync(debugFile[0], "utf-8"));
  const buildInfoName = path.basename(debug.buildInfo);
  const buildInfoPath = await glob(`**/artifacts/build-info/${buildInfoName}`);

  if (buildInfoPath.length !== 1) {
    throw Error(`could not find unique build info for ${contractName}`);
  }

  const possibleSourcePaths = (
    await glob(`**/${contractName}.sol`, {
      ignore: ["artifacts/**", "out/**", "abi/**"],
    })
  ).map((p: string) => {
    if (p.includes("node_modules/")) {
      return p.split("node_modules/")[1];
    } else {
      return p;
    }
  });

  const foundLayouts = [];
  const buildInfo = JSON.parse(fs.readFileSync(buildInfoPath[0], "utf-8"));
  for (const s of possibleSourcePaths) {
    const storageLayout =
      buildInfo.output.contracts[s]?.[contractName]?.storageLayout.storage;
    if (storageLayout) {
      foundLayouts.push(storageLayout);
    }
  }

  if (foundLayouts.length !== 1) {
    throw Error(`could not find unique storage layout for ${contractName}`);
  }

  return foundLayouts[0];
}

async function parsePreDeploy(p: PreDeploy, alloc: any) {
  const execPromise = util.promisify(exec);
  const data = new Map();

  if (p.balance) data.set("balance", p.balance);

  if (p.contract) {
    const artifact = await getArtifact(p.contract);
    data.set("code", artifact.deployedBytecode);
  }

  if (p.storage && p.contract) {
    await getStorageLayout(p.contract);
    const storageLayout = await getStorageLayout(p.contract);

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
