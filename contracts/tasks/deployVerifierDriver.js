const { task } = require("hardhat/config");
const fs = require('fs');
const { exit } = require("process");

const overrides = { gasLimit: 80000000 };

async function deployVerifierDriver() {
    // We get the contract to deploy
    const [account] = await ethers.getSigners();
    const deployerAddress = account.address;
    console.log(`Deploying contracts using ${deployerAddress}`);

    const BlockInitiationVerifier = await ethers.getContractFactory("BlockInitiationVerifier");
    const blockInitiationVerifier = await BlockInitiationVerifier.deploy();
    await blockInitiationVerifier.deployed();

    console.log("BlockInitiationVerifier deployed to:", blockInitiationVerifier.address);

    const BlockFinalizationVerifier = await ethers.getContractFactory("BlockFinalizationVerifier");
    const blockFinalizationVerifier = await BlockFinalizationVerifier.deploy();
    await blockFinalizationVerifier.deployed();

    console.log("BlockFinalizationVerifier deployed to:", blockFinalizationVerifier.address);

    const CallOpVerifier = await ethers.getContractFactory("CallOpVerifier");
    const callOpVerifier = await CallOpVerifier.deploy();
    await callOpVerifier.deployed();

    console.log("CallOpVerifier deployed to:", callOpVerifier.address);

    const EnvironmentalOpVerifier = await ethers.getContractFactory("EnvironmentalOpVerifier");
    const environmentalOpVerifier = await EnvironmentalOpVerifier.deploy();
    await environmentalOpVerifier.deployed();

    console.log("EnvironmentalOpVerifier deployed to:", environmentalOpVerifier.address);

    const InterTxVerifier = await ethers.getContractFactory("InterTxVerifier");
    const interTxVerifier = await InterTxVerifier.deploy();
    await interTxVerifier.deployed();

    console.log("InterTxVerifier deployed to:", interTxVerifier.address);

    const InvalidOpVerifier = await ethers.getContractFactory("InvalidOpVerifier");
    const invalidOpVerifier = await InvalidOpVerifier.deploy();
    await invalidOpVerifier.deployed();

    console.log("InvalidOpVerifier deployed to:", invalidOpVerifier.address);

    const MemoryOpVerifier = await ethers.getContractFactory("MemoryOpVerifier");
    const memoryOpVerifier = await MemoryOpVerifier.deploy(overrides);
    await memoryOpVerifier.deployed();

    console.log("MemoryOpVerifier deployed to:", memoryOpVerifier.address);

    const StackOpVerifier = await ethers.getContractFactory("StackOpVerifier");
    const stackOpVerifier = await StackOpVerifier.deploy(overrides);
    await stackOpVerifier.deployed();

    console.log("StackOpVerifier deployed to:", stackOpVerifier.address);

    const StorageOpVerifier = await ethers.getContractFactory("StorageOpVerifier");
    const storageOpVerifier = await StorageOpVerifier.deploy(overrides);
    await storageOpVerifier.deployed();

    console.log("StorageOpVerifier deployed to:", storageOpVerifier.address);

    const VerifierTestDriver = await ethers.getContractFactory("VerifierTestDriver");
    const verifierTestDriver = await VerifierTestDriver.deploy(
        blockInitiationVerifier.address,
        blockFinalizationVerifier.address,
        interTxVerifier.address,
        stackOpVerifier.address,
        environmentalOpVerifier.address,
        memoryOpVerifier.address,
        storageOpVerifier.address,
        callOpVerifier.address,
        invalidOpVerifier.address,
    );
    await verifierTestDriver.deployed();

    console.log("VerifierTestDriver deployed to:", verifierTestDriver.address);
}

task("deployVerifierDriver", "deploy VerifierTestDriver for testing purpose")
    .setAction(async (taskArgs, { ethers }) => {
        await deployVerifierDriver(ethers);
    });

module.exports = {};