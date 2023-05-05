import { BigNumber, ethers } from "ethers";
import { exec } from "child_process";
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

async function parsePreDeploy(p: PreDeploy, alloc: any) {
  const execPromise = util.promisify(exec);
  const data = new Map();

  if (p.balance) data.set("balance", p.balance);
  if (p.contract) {
    // this allows us to specify any contract within the forge project as pre deploy
    // and makes no assumption about compilation state at the time of genesis file generation
    const { stderr, stdout } = await execPromise(
      `forge inspect ${p.contract} bytes`
    );
    if (stderr || !stdout) Error(stderr);

    data.set("code", stdout.trim());
  }
  if (p.storage && p.contract) {
    const { stderr, stdout } = await execPromise(
      `forge inspect ${p.contract} storage`
    );

    const storageLayout = JSON.parse(stdout);
    if (stderr || !storageLayout) Error(stderr);

    const storage = new Map();
    for (const s of storageLayout.storage) {
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

    // add non variable slots
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
