import dotenv from "dotenv";
import { ethers } from "ethers";
// import Proxy from "../artifacts/@openzeppelin"
import Proxy from "../artifacts/src/pre-deploy/ERC1967Mod.sol/ERC1967Mod.json";
import Placeholder from "../artifacts/src/pre-deploy/UUPSPlaceholder.sol/UUPSPlaceholder.json";
import TinyFaucet from "../artifacts/src/pre-deploy/Faucet.sol/Faucet.json";

dotenv.config({ path: __dirname + "/../.env" });

const main = async () => {
  const proxyAddress = "0xff00000000000000000000000000000000000001";
  const implAddress = "0xff00000000000000000000000000000000000000";
  const provider = new ethers.providers.JsonRpcProvider(
    "http://localhost:4011"
  );
  const sequencer = new ethers.Wallet(
    `0x${process.env.SEQUENCER_PRIVATE_KEY}`,
    provider
  );

  const proxy = new ethers.Contract(proxyAddress, Proxy.abi, provider);

  const impl = new ethers.Contract(implAddress, Placeholder.abi, provider);

  const gasPrice = await provider.getGasPrice();

  const contractBalance = await provider.getBalance(proxyAddress);
  console.log("contract Balance: ", ethers.utils.formatEther(contractBalance));

  console.log(Proxy.abi);

  const r = await proxy.connect(sequencer).getImp({ gasLimit: "300000" });
  console.log({ r });

  // const owner = await impl.connect(proxy).owner();
  // console.log(owner);

  // console.log({ TinyFaucet });
  // contract.connect(sequencer).upgradeTo(TinyFaucet.bytecode);

  // const allowedAmount = await contract.amountAllowed();
  // console.log("amount allowed: ", allowedAmount);

  // // Transfers amountAllowed from faucet to signer address
  // const requestFundTx = await contract
  //   .connect(signer)
  //   .requestFunds(signer.address, {
  //     gasLimit: "300000",
  //   });
  // await requestFundTx.wait();

  // // checks balance after requestFunds request
  // signerbalance = await provider.getBalance(signer.address);
  // console.log("sequencer Balance: ", ethers.utils.formatEther(signerbalance));

  // contractBalance = await provider.getBalance(contractAddress);
  // console.log("contract Balance: ", ethers.utils.formatEther(contractBalance));

  // // Transfers faucet balance to sequencer
  // const retrieve = await contract.connect(signer).retrieve({
  //   gasLimit: "260000",
  //   gasPrice,
  // });
  // await retrieve.wait();

  // // checking balance after retrieving complete balance from faucet
  // signerbalance = await provider.getBalance(signer.address);
  // console.log("sequencer Balance: ", ethers.utils.formatEther(signerbalance));

  // contractBalance = await provider.getBalance(contractAddress);
  // console.log("contract Balance: ", ethers.utils.formatEther(contractBalance));
};

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
