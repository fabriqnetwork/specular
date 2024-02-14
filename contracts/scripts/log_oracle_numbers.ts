import { ethers } from "hardhat";
import { program } from "commander";

const L1_ORACLE_ADDRESS = "0x2A00000000000000000000000000000000000010";

program
  .option("--l1-rpc <string>")
  .option("--l2-rpc <string>")
  .option("--interval <number>");

async function main() {
  program.parse(process.argv);
  const options = program.opts();

  options.l1Rpc = options.l1Rpc || "http://127.0.0.1:8545";
  options.l2Rpc = options.l2Rpc || "http://127.0.0.1:4011";
  options.interval = options.interval || 2000;

  const providers = {
    l1: new ethers.providers.JsonRpcProvider(options.l1Rpc),
    l2: new ethers.providers.JsonRpcProvider(options.l2Rpc),
  };

  const l1OracleFactory = await ethers.getContractFactory("L1Oracle");
  const l1Oracle = l1OracleFactory
    .attach(L1_ORACLE_ADDRESS)
    .connect(providers.l2);

  while (true) {
    const log = {
      time: Date.now(),
      l1BlockLatest: await providers.l1.getBlockNumber(),
      l2BlockLatest: await providers.l2.getBlockNumber(),
      l1OracleBlock: (await l1Oracle.number()).toNumber(),
    };

    console.log(JSON.stringify(log));

    await new Promise((r) => setTimeout(r, options.interval));
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
