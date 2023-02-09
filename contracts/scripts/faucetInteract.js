const ethers = require('ethers');
const TinyFaucet = require('../artifacts/src/pre-deploy/Faucet.sol/Faucet.json');

const main = async() => {

  const contractAddress = "0x0000000000000000000000000000000000000020"
  const provider = new ethers.providers.JsonRpcProvider("http://localhost:4011");

  const signer = new ethers.Wallet("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",provider);
  const contract = new ethers.Contract(contractAddress, TinyFaucet.abi, provider);
  const gasPrice = await provider.getGasPrice();
  
  // initial balance of sequencer & faucet
  let signerbalance = await provider.getBalance(signer.address);
  console.log("sequencer Balance: ",ethers.utils.formatEther(signerbalance))

  contractBalance = await provider.getBalance(contractAddress);
  console.log("contract Balance: ",ethers.utils.formatEther(contractBalance))

  const owner = await contract.owner();
  console.log(owner)

  const allowedAmount = await contract.amountAllowed();
  console.log("amount allowed: ",allowedAmount)
  
  // Transfers amountAllowed from faucet to signer address
  const requestFundTx = await contract.connect(signer).requestFunds(signer.address, {
    gasLimit: "300000",
  });
  await requestFundTx.wait();
  
  // checks balance after requestFunds request
  signerbalance = await provider.getBalance(signer.address);
  console.log("sequencer Balance: ",ethers.utils.formatEther(signerbalance))
  
  contractBalance = await provider.getBalance(contractAddress);
  console.log("contract Balance: ",ethers.utils.formatEther(contractBalance))

  // Transfers faucet balance to sequencer 
  let retrieve = await contract.connect(signer).retrieve({
    gasLimit: "260000",
    gasPrice,
  });
  await retrieve.wait(); 

  // checking balance after retrieving complete balance from faucet
  signerbalance = await provider.getBalance(signer.address);
  console.log("sequencer Balance: ",ethers.utils.formatEther(signerbalance))

  contractBalance = await provider.getBalance(contractAddress);
  console.log("contract Balance: ",ethers.utils.formatEther(contractBalance))


}


main();