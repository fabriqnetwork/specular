// rules which govern clef key management behavior
// for more info, see: https://geth.ethereum.org/docs/tools/clef/rules

function ApproveTx(r) {
  // all transactions not going to the rollup addr are rejected
  // TODO: do we need additional / more fine grained rules?

  // this file can't access env variables so addresses have to be hardcoded
  // if (r.transaction.to.toLowerCase() == '0x5fc8d32690cc91d4c39d9d3abcbd16989f875707') {
  //   return 'Approve';
  // }

  return 'Approve'
}

function ApproveListing() {
  return 'Approve';
}

function ApproveSignData() {
  return 'Approve';
}
