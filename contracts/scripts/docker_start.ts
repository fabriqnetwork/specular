import { executeCommand } from "./docker_utils";
import path from "path";

const ROOT_DIR = path.join(__dirname, "../../");
const SPECULAR_DATADIR =
  process.env.SPECULAR_DATADIR || path.join(ROOT_DIR, "specular-datadir");

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
