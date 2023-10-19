import { BigNumber } from "ethers";

export class RelayerState {
  // For throttling the update frequency of the L1 oracle
  lastSentL1OracleBlockNumber: BigNumber = BigNumber.from(0);
  // The last block number that the L1 oracle was updated to
  lastUpdatedL1OracleBlockNumber: BigNumber = BigNumber.from(0);

  constructor() {}

  sentL1OracleValues(blockNumber: BigNumber) {
    this.lastSentL1OracleBlockNumber = blockNumber;
  }

  updatedL1OracleValues(blockNumber: BigNumber) {
    this.lastUpdatedL1OracleBlockNumber = blockNumber;
  }
}
