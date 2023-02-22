import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import { Manifest } from "@openzeppelin/upgrades-core";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, getNamedAccounts, ethers, upgrades, network } = hre;
  const { save } = deployments;
  const { sequencer, deployer } = await getNamedAccounts();
  const deployerSigner = await ethers.getSigner(deployer);
  const { provider } = network;

  const sequencerInboxProxyAddress = (await deployments.get("SequencerInbox"))
    .address;
  const verifierProxyAddress = (await deployments.get("Verifier")).address;

  const rollupArgs = [
    sequencer, // address _vault
    sequencerInboxProxyAddress, // address _sequencerInbox
    verifierProxyAddress, // address _verifier
    "0x0000000000000000000000000000000000000000", // address _stakeToken
    5, // uint256 _confirmationPeriod
    0, // uint256 _challengePeriod
    0, // uint256 _minimumAssertionPeriod
    1000000000000, // uint256 _maxGasPerAssertion
    0, // uint256 _baseStakeAmount
    "0x744c19d2e8593c97867b3b6a3588f51cd9dbc5010a395cf199be4bbb353848b8", // bytes32 _initialVMhash
  ];

  const Rollup = await ethers.getContractFactory("Rollup", deployer);
  const rollup = await upgrades.deployProxy(Rollup, rollupArgs, {
    initializer: "initialize",
    from: sequencer,
    timeout: 0,
    kind: "uups",
  });

  await rollup.deployed();
  console.log("Rollup Proxy:", rollup.address);
  console.log(
    "Rollup Implementation Address",
    await upgrades.erc1967.getImplementationAddress(rollup.address)
  );
  console.log(
    "Rollup Admin Address",
    await upgrades.erc1967.getAdminAddress(rollup.address)
  );

  const artifact = await deployments.getExtendedArtifact("Rollup");
  const proxyDeployments = {
    address: rollup.address,
    ...artifact,
  };
  await save("Rollup", proxyDeployments);
};

export default func;
func.tags = ["Rollup"];
