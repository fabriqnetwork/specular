import { exec } from 'child_process';
import path from 'path';

const ROOT_DIR = path.join(__dirname, "../../");
const SPECULAR_DATADIR =
  process.env.SPECULAR_DATADIR || path.join(ROOT_DIR, "specular-datadir");

async function executeCommand(command: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const child = exec(command, (error, stdout, stderr) => {
      if (error) {
        console.error(`Error executing command: ${error.message}`);
        reject(error);
        return;
      }
      resolve();
    });

    // Log stdout
    child.stdout.on('data', (data) => {
      console.log(`stdout: ${data}`);
    });

    // Log stderr
    child.stderr.on('data', (data) => {
      console.error(`stderr: ${data}`);
    });
  });
}

// Start containers
async function startContainers() {
  console.log("In startContainers(), this is ROOT_DIR: ", ROOT_DIR);
  console.log(
    "In startContainers(), this is SPECULAR_DATADIR: ",
    SPECULAR_DATADIR
  );
  const command = `docker compose -f ${ROOT_DIR}docker/docker-compose-integration-test.yml up -d --build`;
  await executeCommand(command);
  console.log("Done with startContainers()..");
}

startContainers();
