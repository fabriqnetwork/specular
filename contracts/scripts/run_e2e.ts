import { exec } from "child_process";
import * as fs from "fs";
import { promisify } from "util";

const PROJECT_PATH = "project";
const CONTRACTS_PATH = "contracts";

// For err handling
const execAsync = promisify(exec);

// Check docker version
async function checkDockerVersion(): Promise<void> {
  console.log("Running checkDockerVersion()");
  const { stdout, stderr } = await execAsync("docker version");
  console.log("checkDockerVersion stdout: ", stdout);
  if (stderr) {
    console.error("checkDockerVersion stderr:", stderr);
  }
}

// Install docker compose
async function installDockerCompose(version: string): Promise<void> {
  console.log("Running installDockerCompose()");
  await execAsync(
    `sudo curl -L "https://github.com/docker/compose/releases/download/${version}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose`
  );
  await execAsync("sudo chmod +x /usr/local/bin/docker-compose");
}

// Update git submodules
async function updateGitSubmodules() {
  console.log("Running updateGitSubmodules()");
  return new Promise<void>((resolve, reject) => {
    exec(
      "git submodule init && git submodule update --recursive",
      (error, stdout, stderr) => {
        if (error) {
          console.log("updateGitSubmodules stderr: ", stderr);
          reject(error);
        } else {
          console.log("updateGitSubmodules stdout: ", stdout);
          resolve();
        }
      }
    );
  });
}

// Checkout code
async function checkoutCode(): Promise<void> {
  console.log("Running checkoutCode()");
  await execAsync("git submodule init");
  await execAsync("git submodule update --recursive");
}

// List files in current dir
async function listFilesInCurrentDir(): Promise<void> {
  console.log("Running listFilesInCurrentDir()");
  const { stdout } = await execAsync("pwd && ls -la");
  console.log(stdout);
}

// Start up the containers
async function startupContainers(): Promise<void> {
  console.log("Running startupContainers()");
  const { stdout: pwdOutput } = await execAsync("pwd && ls -la");
  console.log(pwdOutput);
  await execAsync("cd project");
  await execAsync("rm -rf project/specular-datadir/geth/");
  await execAsync("npx ts-node ../contracts/scripts/docker_start.ts");
  await execAsync(
    "wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh"
  );
  await execAsync("chmod +x wait-for-it.sh");
  await execAsync("sleep 30");
  await execAsync("docker ps");
  const { stdout: hardhatLogs } = await execAsync("docker logs hardhat");
  console.log(hardhatLogs);
  const { stdout: sequencerLogs } = await execAsync("docker logs sequencer");
  console.log(sequencerLogs);
  await execAsync("./wait-for-it.sh -t 240 127.0.0.1:8545");
  await execAsync("./wait-for-it.sh -t 240 127.0.0.1:4011");
  const { stdout: hardhatLogs2 } = await execAsync("docker logs hardhat");
  console.log(hardhatLogs2);
  const { stdout: sequencerLogs2 } = await execAsync("docker logs sequencer");
  console.log(sequencerLogs2);
}

// Run testing script
async function runTestingScript(): Promise<void> {
  console.log("Running runTestingScript()");

  await exec("yarn install", { cwd: CONTRACTS_PATH });

  await exec("npx hardhat deploy --network localhost", {
    cwd: CONTRACTS_PATH,
  });

  await exec(
    "yarn run ts-node /home/runner/work/specular/specular/contracts/scripts/testing.ts",
    { cwd: CONTRACTS_PATH }
  );

  await exec("docker logs hardhat");
  await exec("docker logs sequencer");
}

// Stop and remove containers
async function stopAndRemoveContainers(): Promise<void> {
  console.log("Running stopAndRemoveContainers()");

  await exec("npx ts-node contracts/scripts/docker_start.ts --stop --remove", {
    cwd: PROJECT_PATH,
  });
}

// Run e2e process
async function runE2E(): Promise<void> {
  await checkDockerVersion();
  await installDockerCompose("1.29.2");
  await checkoutCode();
  await updateGitSubmodules();
  await listFilesInCurrentDir();
  await startupContainers();
  await runTestingScript();
  await stopAndRemoveContainers();
}

runE2E().catch((error) => {
  console.error(error);
  process.exit(1);
});
