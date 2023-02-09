const ethers = require('ethers');
const fs = require('fs');
const faucet = require("../artifacts/src/pre-deploy/Faucet.sol/Faucet.json");
const genesisPath = "../../clients/geth/specular/data/genesis.json";
const genesisJson = require(genesisPath);
const assert = require("assert");

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
    const faucetDeployedBytecode = faucet.deployedBytecode;
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
    const faucetAddress = "0x0000000000000000000000000000000000000020";
    
    genesisJson.alloc[faucetAddress.toString()] = createFaucetContractObject();
    fs.writeFileSync(genesisPath, JSON.stringify(genesisJson, null, 2));

}

main();