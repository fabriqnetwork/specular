const { task } = require("hardhat/config");
const fs = require('fs');
const { exit } = require("process");

const overrides = { gasLimit: 8000000 };

const emtpyHash = "0x0000000000000000000000000000000000000000000000000000000000000000";

async function main(ethers, addr, proofFile) {
    const verifierTestDriver = await ethers.getContractAt("VerifierTestDriver", addr);
    
    const { ctx, proof } = JSON.parse(fs.readFileSync(proofFile));
    console.log(`processing tx ${ctx.txnHash}`);

    const transaction = [
        ctx.txNonce,
        ctx.gasPrice,
        ctx.gas,
        ctx.recipient,
        ctx.value,
        ctx.input,
        ctx.txV,
        ctx.txR,
        ctx.txS,
    ];
    
    const res = await verifierTestDriver.verifyProof(
        ctx.coinbase,
        ctx.timestamp,
        ctx.blockNumber,
        ctx.origin,
        ctx.txnHash,
        transaction,
        proof.verifier,
        proof.currHash,
        proof.proof,
    );

    console.log(await res.wait());
}

task("verifyOsp", "verify osp")
    .addParam("addr", "VerifierTestDriver contract address")
    .addParam("proof", "Path to proof file")
    .setAction(async (taskArgs, { ethers }) => {
        await main(ethers, taskArgs.addr, taskArgs.proof);
    });

module.exports = {};