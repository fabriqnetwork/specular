import { executeCommand } from "./docker_utils";

const ROOT_DIR = __dirname + "/../../";

// Start containers
async function startContainers() {
  const command = `docker-compose -f ${ROOT_DIR}docker/docker-compose-integration-test.yml up -d --build`;
  await executeCommand(command);
}

startContainers();
