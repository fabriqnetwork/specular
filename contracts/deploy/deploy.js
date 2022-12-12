module.exports = async (hre) => {
  const { getNamedAccounts, ethers, upgrades } = hre
  const { sequencer } = await getNamedAccounts()

  const Inbox = await ethers.getContractFactory('SequencerInbox')
  const inbox = await upgrades.deployProxy(Inbox, [sequencer], { initializer: 'initialize', from: sequencer })

  console.log('Sequencer Inbox Proxy:', inbox.address)
  console.log('Sequencer Inbox Implementation Address', await upgrades.erc1967.getImplementationAddress(inbox.address))
  console.log('Sequencer Inbox Admin Address', await upgrades.erc1967.getAdminAddress(inbox.address))

  const Verifier = await ethers.getContractFactory('Verifier')
  const verifier = await upgrades.deployProxy(Verifier, [], { initializer: 'initialize' })

  console.log('Verifier Proxy:', verifier.address)
  console.log('Verifier Implementation Address', await upgrades.erc1967.getImplementationAddress(verifier.address))
  console.log('Verifier Admin Address', await upgrades.erc1967.getAdminAddress(verifier.address))

  const rollupArgs = [
    sequencer, // address _owner
    inbox.address, // address _sequencerInbox,
    verifier.address, // address _verifier,
    '0x0000000000000000000000000000000000000000', // address _stakeToken,
    5, // uint256 _confirmationPeriod,
    0, // uint256 _challengePeriod,
    0, // uint256 _minimumAssertionPeriod,
    1000000000000, // uint256 _maxGasPerAssertion,
    0, // uint256 _baseStakeAmount
    '0x744c19d2e8593c97867b3b6a3588f51cd9dbc5010a395cf199be4bbb353848b8' // bytes32 _initialVMhash
  ]

  const Rollup = await ethers.getContractFactory('Rollup')
  const rollup = await upgrades.deployProxy(Rollup, rollupArgs, { initializer: 'initialize' })

  console.log('Rollup Proxy:', rollup.address)
  console.log('Rollup Implementation Address', await upgrades.erc1967.getImplementationAddress(rollup.address))
  console.log('Rollup Admin Address', await upgrades.erc1967.getAdminAddress(rollup.address))
}

module.exports.tags = ['SequencerInbox']
