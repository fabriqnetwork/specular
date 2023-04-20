import { executeCommand } from "./docker_utils";
import path from "path";

const ROOT_DIR = path.join(__dirname, "../../");

// Stop containers
async function stopContainers() {
  console.log("In stopContainers(), this is ROOT_DIR: ", ROOT_DIR);
  const command = `docker compose -f ${ROOT_DIR}docker/docker-compose-integration-test.yml down`;
  await executeCommand(command);
  console.log("done with stopContainers()..");
}

stopContainers();
