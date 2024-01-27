import fs from "fs";
import { parseFlag } from "./utils";

require("dotenv").config();

// TODO(#304): consider moving to golang (ops).
async function main() {
  const deploymentsConfig = parseFlag("--deployments-config-path");
  const deploymentsPath = parseFlag("--deployments", "./deployments/localhost");
  await generateContractAddresses(deploymentsConfig, deploymentsPath);
}

/**
 * Reads the L1 deployment and writes deployments address to the deployments env file
 */
export async function generateContractAddresses(
  deploymentsConfigPath: string,
  deploymentsPath: string,
) {
  // check the deployments dir - error out if it is not there
  const deploymentFiles = fs.readdirSync(deploymentsPath);
  let result = "";
  for (const deploymentFile of deploymentFiles) {
    if (
      deploymentFile.startsWith("Proxy__") &&
      deploymentFile.endsWith(".json")
    ) {
      const deployment = JSON.parse(
        fs.readFileSync(`${deploymentsPath}/${deploymentFile}`, "utf-8"),
      );
      let contractName = deploymentFile
        .replace(/^Proxy__/, "")
        .replace(/\.json$/, "");
      contractName = contractName
        .replace(/([a-z])([A-Z])/g, "$1_$2")
        .toUpperCase();
      result += `${contractName}_ADDR=${deployment.address}\n`;
    }
  }
  fs.writeFileSync(deploymentsConfigPath, result);
  console.log(`successfully wrote deployments to: ${deploymentsConfigPath}`);
}

if (!require.main!.loaded) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
