const ethers = require('ethers');
const storageContract = require('./artifacts/src/storage.sol/Storage.json');

const main = async() => {

  const contractAddress = "0x0000000000000000000000000000000000000020"
  const provider = new ethers.providers.JsonRpcProvider("http://localhost:4011");
  const code = await provider.getCode(contractAddress)
  const signer = new ethers.Wallet("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",provider);
  const contract = new ethers.Contract(contractAddress, storageContract.abi, provider);

  let contractBalance = await provider.getBalance(contractAddress);
  console.log("contract Balance: ",ethers.utils.formatEther(contractBalance))
  // console.log(contract)

  let retrieve = await contract.retrieve();
  console.log("Number retrieved",Number(retrieve))
  const numStore = await contract.connect(signer).store(43);
  await numStore.wait();
  retrieve = await contract.retrieve();
  console.log("Number retrieved",Number(retrieve))

  const getPaid = await contract.connect(signer).getPaid();
  await getPaid.wait();
  contractBalance = await provider.getBalance(contractAddress);
  console.log("contract Balance: ",ethers.utils.formatEther(contractBalance))
  const signerbalance = await provider.getBalance(signer.address);
  console.log("sequencer Balance: ",ethers.utils.formatEther(signerbalance))
}


main();