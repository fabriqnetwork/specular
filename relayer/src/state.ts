import { CrossDomainMessage } from "./messager";
import { BigNumber } from "ethers";
import { PriorityQueue } from "typescript-collections";

type PendingDeposit = {
  l1BlockNumber: number;
  depositHash: string;
  depositTx: CrossDomainMessage;
};

type PendingWithdrawal = {
  l2BlockNumber: number;
  withdrawalHash: string;
  withdrawalTx: CrossDomainMessage;
};

function depositPriority(a: PendingDeposit, b: PendingDeposit): number {
  // smaller block number gets higher priority
  return b.l1BlockNumber - a.l1BlockNumber;
}

function withdrawalPriority(
  a: PendingWithdrawal,
  b: PendingWithdrawal
): number {
  // smaller block number gets higher priority
  return b.l2BlockNumber - a.l2BlockNumber;
}

type CachedDepositProof = {
  depositAccountProof: string[];
};

type CachedWithdrawalProof = {
  encodedBlockHeader: string;
  withdrawalAccountProof: string[];
};

type L2BlockNumberMappingEntry = {
  l2BlockNumber: number;
  inboxSize: BigNumber;
};

function l2BlockNumberMappingEntryPriority(
  a: L2BlockNumberMappingEntry,
  b: L2BlockNumberMappingEntry
): number {
  // smaller inbox size gets higher priority
  const diff = a.inboxSize.sub(b.inboxSize);
  if (diff.isZero()) {
    return 0;
  }
  if (diff.isNegative()) {
    return 1;
  }
  return -1;
}

export class RelayerState {
  pendingDeposits: PriorityQueue<PendingDeposit> =
    new PriorityQueue<PendingDeposit>(depositPriority);
  pendingWithdrawals: PriorityQueue<PendingWithdrawal> =
    new PriorityQueue<PendingWithdrawal>(withdrawalPriority);

  // For throttling the update frequency of the L1 oracle
  lastSentL1OracleBlockNumber: number = 0;
  // The last block number that the L1 oracle was updated to
  lastUpdatedL1OracleBlockNumber: number = 0;

  l2BlockNumberMapping: PriorityQueue<L2BlockNumberMappingEntry> =
    new PriorityQueue<L2BlockNumberMappingEntry>(
      l2BlockNumberMappingEntryPriority
    );
  lastConfirmedAssertionID: BigNumber = BigNumber.from(0);
  lastConfirmedL2BlockNumber: number = 0;
  lastConfirmedInboxSize: BigNumber = BigNumber.from(0);
  lastConfirmedVMHash: string = "";
  lastConfirmedL2GasUsed: BigNumber = BigNumber.from(0);

  assertionMap: Map<string, [BigNumber, string]> = new Map<
    string,
    [BigNumber, string]
  >();

  constructor() {}

  updateL2BlockNumberMapping(l2BlockNumber: number, inboxSize: BigNumber) {
    console.log(
      "see inbox event, mapping l2 block number",
      l2BlockNumber,
      "to inbox size",
      inboxSize.toString()
    );
    this.l2BlockNumberMapping.add({
      l2BlockNumber,
      inboxSize,
    });
  }

  updateCreatedAssertion(
    assertionID: BigNumber,
    l2GasUsed: BigNumber,
    vmHash: string
  ) {
    console.log("see assertion created event", assertionID.toString());
    this.assertionMap.set(assertionID.toString(), [l2GasUsed, vmHash]);
  }

  updateConfirmedInboxSize(assertionID: BigNumber, inboxSize: BigNumber) {
    console.log("see assertion confirmed event", assertionID.toString());
    this.lastConfirmedAssertionID = assertionID;
    this.lastConfirmedInboxSize = inboxSize;
    try {
    [this.lastConfirmedL2GasUsed, this.lastConfirmedVMHash] =
      this.assertionMap.get(assertionID.toString())!;
    } catch (e) {
      console.error(e)
      console.log("assertionID ", assertionID.toString(), " not found, skipping")
      return;
    }
    this.assertionMap.delete(assertionID.toString());
    while (!this.l2BlockNumberMapping.isEmpty()) {
      const entry = this.l2BlockNumberMapping.peek()!;
      if (entry.inboxSize.lt(inboxSize)) {
        this.lastConfirmedL2BlockNumber = entry.l2BlockNumber;
        this.l2BlockNumberMapping.dequeue();
      } else {
        break;
      }
    }
    if (this.l2BlockNumberMapping.isEmpty()) {
      console.error("dubious error: confirmed a non-existing batch");
      return;
    }
    const entry = this.l2BlockNumberMapping.dequeue()!;
    if (!entry.inboxSize.eq(inboxSize)) {
      console.log(
        "dubious error: assertion created in the middle of the batch"
      );
      console.log(
        "confirmed assertionID",
        assertionID.toString(),
        "confirmed inbox size",
        inboxSize.toString(),
        "entry",
        entry
      );
      return;
    }
    this.lastConfirmedL2BlockNumber = entry.l2BlockNumber;
  }

  sentL1OracleValues(blockNumber: number) {
    this.lastSentL1OracleBlockNumber = blockNumber;
  }

  updatedL1OracleValues(blockNumber: number) {
    this.lastUpdatedL1OracleBlockNumber = blockNumber;
  }

  addDeposit(
    l1BlockNumber: number,
    depositHash: string,
    depositTx: CrossDomainMessage
  ) {
    this.pendingDeposits.add({
      l1BlockNumber,
      depositHash,
      depositTx,
    });
  }

  readdDeposit(deposit: PendingDeposit) {
    this.pendingDeposits.add(deposit);
  }

  getNextDepositBlockNumber(): number | undefined {
    return this.pendingDeposits.peek()?.l1BlockNumber;
  }

  getNextDeposit(): PendingDeposit | undefined {
    return this.pendingDeposits.dequeue();
  }

  addWithdrawal(
    l2BlockNumber: number,
    withdrawalHash: string,
    withdrawalTx: CrossDomainMessage
  ) {
    this.pendingWithdrawals.add({
      l2BlockNumber,
      withdrawalHash,
      withdrawalTx,
    });
  }

  readdWithdrawal(withdrawal: PendingWithdrawal) {
    this.pendingWithdrawals.add(withdrawal);
  }

  getNextWithdrawalBlockNumber(): number | undefined {
    return this.pendingWithdrawals.peek()?.l2BlockNumber;
  }

  getNextWithdrawal(): PendingWithdrawal | undefined {
    return this.pendingWithdrawals.dequeue();
  }
}
