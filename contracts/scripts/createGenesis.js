const ethers = require('ethers');
const fs = require('fs');
const FaucetJson = require("../artifacts/src/pre-deploy/Faucet.sol/Faucet.json");
const assert = require("assert");
let GenesisJson;

const createContractObject = (deployedBytecode, contractBalance, storageSlots, valueAtSlots) => {
    
    assert(storageSlots.length == valueAtSlots.length, "incorrect storage-values array lengths")
    
    let storageSlotsObj = {};
    for(let i = 0; i < storageSlots.length; i++) {
        storageSlotsObj[storageSlots[i].toString()] = valueAtSlots[i];
    }

    const contractObject = {
        "code": deployedBytecode,
        "balance": Number(contractBalance).toString(),
        "storage": storageSlotsObj
    }
    return contractObject;
}

const createFaucetContractObject = () => {
    const faucetDeployedBytecode = FaucetJson.deployedBytecode;
    const faucetBalance = ethers.BigNumber.from("10").pow(20);
    
    let storageSlots = [];
    let valueAtSlots = [];

    storageSlots[0] = "0x0000000000000000000000000000000000000000000000000000000000000000";
    storageSlots[1] = "0x0000000000000000000000000000000000000000000000000000000000000001";
    valueAtSlots[0] = "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266";
    valueAtSlots[1] = "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000"
    
    const faucetObject = createContractObject(faucetDeployedBytecode, faucetBalance, storageSlots, valueAtSlots);
    console.log(faucetObject)
    return faucetObject;
}

const main = () => {

    const inFlagIndex = process.argv.indexOf("--in");
    let baseGenesisPath;
    
    if(inFlagIndex > -1) {
        baseGenesisPath = process.argv[inFlagIndex+1];
        GenesisJson = require(baseGenesisPath);
    } else {
        throw new Error("Please specify the base genesis path");
    }

    const outFlagIndex = process.argv.indexOf("--out");
    let genesisPath;

    if(outFlagIndex > -1) {
        genesisPath = process.argv[outFlagIndex+1];
    } else {
        console.log("Setting out genesis path same as base genesis path");
        genesisPath = baseGenesisPath;
    }
    
    const faucetAddress = "0x0000000000000000000000000000000000000020";
    
    GenesisJson.alloc[faucetAddress.toString()] = createFaucetContractObject();
    fs.writeFileSync(genesisPath, JSON.stringify(GenesisJson, null, 2));

}

main();