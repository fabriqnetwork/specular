import { executeCommand } from "./docker_utils";

const ROOT_DIR = __dirname + "/../../";

// Stop containers
async function stopContainers() {
  const command = `docker-compose -f ${ROOT_DIR}docker/docker-compose-integration-test.yml down`;
  await executeCommand(command);
}

stopContainers();
