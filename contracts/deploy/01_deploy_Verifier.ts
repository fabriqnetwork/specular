import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployments, getNamedAccounts, ethers, upgrades} = hre;
    const {deploy} = deployments;
    const {sequencer} = await getNamedAccounts();
    
    console.log(sequencer);

    const Verifier = await ethers.getContractFactory("VerifierEntry");
    const verifier = await upgrades.deployProxy(Verifier, [], {initializer: 'initialize', from: sequencer, timeout: 0});
    
    await verifier.deployed();
    console.log("Verifier Proxy:", verifier.address);
    console.log("Verifier Implementation Address", await upgrades.erc1967.getImplementationAddress(verifier.address));
    console.log("Verifier Admin Address", await upgrades.erc1967.getAdminAddress(verifier.address))    
  
}

export default func;
func.tags = ['Verifier'];