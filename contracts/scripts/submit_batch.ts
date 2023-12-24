const ethers = require("ethers");

const sendBatch = async () => {
  const provider = new ethers.providers.JsonRpcBatchProvider(
    "http://localhost:4011"
  );

  const priv =
    "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80";
  const addr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";

  const signer = new ethers.Wallet(priv, provider);
  const balance = await provider.getBalance(addr);
  console.log("balance: " + balance.toString());

  let nonce = await provider.getTransactionCount(signer.address);
  nonce += 0;
  let promises = [];
  for (let i = 0; i < 10; i++) {
    promises.push(
      signer.sendTransaction({
        to: "0xf112347faDA222A95d84626b19b2af1DB6672C18",
        value: ethers.utils.parseEther("0.1"),
        nonce: nonce,
      })
    );
    nonce++;
  }

  const results = await Promise.all(promises);
  console.log(results);
};

sendBatch();
