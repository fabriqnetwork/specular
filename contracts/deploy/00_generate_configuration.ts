import { generateGenesisFile } from "../scripts/create_genesis.ts";
import { exec } from "child_process";
import util from "node:util";
import path from "path";
import fs from "fs";

const CLIENT_SBIN_DIR = `${__dirname}/../../clients/geth/specular/sbin`;
const CLIENT_DATA_DIR = `${__dirname}/../../clients/geth/specular/data`;
const execPromise = util.promisify(exec);

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  // 1. generate genesis file
  const dataPath = `${__dirname}/../../clients/geth/specular/data/`;
  const baseGenesisPath = path.join(dataPath, "base_genesis.json");
  const genesisPath = path.join(dataPath, "genesis.json");

  // 2. init geth nodes
  await execPromise(path.join(CLIENT_SBIN_DIR, "clean.sh"));
  await execPromise(path.join(CLIENT_SBIN_DIR, "init.sh"));

  // 3. generate genesis hash
  const { err, stdout } = await execPromise(
    path.join(CLIENT_SBIN_DIR, "export_genesis.sh")
  );
  const initialVmHash = JSON.parse(stdout).root;

  if (err !== undefined || !initialVmHash) {
    throw Error("could not export genesis hash", err);
  }

  // 4. write genesis hash to file for use in later deployments
  fs.writeFileSync(
    path.join(CLIENT_DATA_DIR, "initialVmHash.txt"),
    initialVmHash,
    "utf-8"
  );
};

export default func;
func.tags = ["GenerateConfig"];
