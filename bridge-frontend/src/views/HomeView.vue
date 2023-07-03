<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import type { Ref } from "vue";
import type {
  PendingDeposit,
  PendingWithdrawal,
  MessageProof,
} from "@/client/types";
import type { BigNumber } from "ethers";
import type { JsonRpcProvider } from "@ethersproject/providers";
import {
  useAccount,
  useConnect,
  useEnsName,
  useSwitchNetwork,
  useNetwork,
  useProvider,
} from "vagmi";
import { InjectedConnector } from "vagmi/connectors/injected";
import { chiadoChain, specularChain } from "@/client/chains";
import { ethers } from "ethers";
import {
  L1PORTAL_ADDRESS,
  L2PORTAL_ADDRESS,
  ROLLUP_ADDRESS,
  L1ORACLE_ADDRESS,
  INBOX_ADDRESS,
  DEPOSIT_BALANCE_THRESHOLD,
} from "@/client/constants";
import { transactor } from "@/client/providers";
import { getStorageKey, requestFundDeposit } from "@/client/utils";
import {
  IL1Portal__factory,
  IL2Portal__factory,
  L1Oracle__factory,
  IRollup__factory,
  ISequencerInbox__factory,
} from "../../../contracts/typechain-types";

const { address, isConnected } = useAccount();
const { data: ensName } = useEnsName({
  address,
});
const { connect } = useConnect({
  connector: new InjectedConnector(),
});
const { chains, error, isLoading, pendingChainId, switchNetworkAsync } =
  useSwitchNetwork();
const { chain } = useNetwork();

const l1Provider = useProvider({ chainId: chiadoChain.id });
const l2Provider = useProvider({ chainId: specularChain.id });

const mode = ref("Deposit");
const currentChain = computed(() => {
  if (mode.value === "Deposit") {
    return chiadoChain;
  } else {
    return specularChain;
  }
});
const targetChain = computed(() => {
  if (mode.value === "Deposit") {
    return specularChain;
  } else {
    return chiadoChain;
  }
});
const amount = ref(0);
const l1Balance = ref(0);
const l2Balance = ref(0);

const currentChainBalance = computed(() => {
  if (mode.value === "Deposit") {
    return l1Balance.value;
  } else {
    return l2Balance.value;
  }
});
const targetChainBalance = computed(() => {
  if (mode.value === "Deposit") {
    return l2Balance.value;
  } else {
    return l1Balance.value;
  }
});
const isInsufficientBalance = computed(
  () => currentChainBalance.value < Number(amount.value)
);
const bridgeBtnText = computed(() =>
  isInsufficientBalance.value
    ? "Insufficient Balance"
    : isZeroAmount.value
    ? "Zero Amount"
    : isNaNAmount.value
    ? "Not a Number"
    : mode.value
);
const isZeroAmount = computed(() => Number(amount.value) === 0);
const isNaNAmount = computed(() => Number.isNaN(Number(amount.value)));

const depositTxHash = ref("");
const withdrawTxHash = ref("");
const pendingDeposit: Ref<PendingDeposit | undefined> = ref(undefined);
const pendingWithdrawal: Ref<PendingWithdrawal | undefined> = ref(undefined);
const depositFinalizedTxHash = ref("");
const withdrawFinalizedTxHash = ref("");
const isDepositOngoing = ref(false);
const isWithdrawOngoing = ref(false);
const isDepositReadyToFinalize = ref(false);
const isWithdrawalReadyToFinalize = ref(false);
const isDepositFinalizationRequested = ref(false);
const isWithdrawalFinalizationRequested = ref(false);
const isDepositFinished = ref(false);
const isWithdrawFinished = ref(false);

const isDepositError = ref(false);
const isWithdrawalError = ref(false);

const sequencerInboxInterface = ISequencerInbox__factory.createInterface();
const inboxSizeToBlockNumberMap = new Map<string, number>();

async function refreshL1Balance() {
  if (address.value) {
    const balance = await l1Provider.value.getBalance(address.value);
    l1Balance.value = Number(
      ethers.utils.formatUnits(
        balance,
        chiadoChain.nativeCurrency?.decimals
      )
    );
  }
}

async function refreshL2Balance() {
  if (address.value) {
    const balance = await l2Provider.value.getBalance(address.value);
    l2Balance.value = Number(
      ethers.utils.formatUnits(
        balance,
        specularChain.nativeCurrency?.decimals
      )
    );
  }
}

async function refreshBalance() {
  try {
    await refreshL1Balance();
  } catch (e) {
    console.error(e);
  }
  try {
    await refreshL2Balance();
  } catch (e) {
    console.error(e);
  }
  console.log(currentChainBalance.value, targetChainBalance.value);
}

async function bridge() {
  if (!isConnected.value) {
    await connect.value();
  }
  if (Number.isNaN(Number(amount.value))) {
    return;
  }
  if (isInsufficientBalance.value) {
    return;
  }
  if (mode.value === "Deposit") {
    const signer = transactor.getSigner();
    const tx = await signer.sendTransaction({
      to: L1PORTAL_ADDRESS,
      value: ethers.utils.parseUnits(
        Number(amount.value).toString(),
        currentChain.value.nativeCurrency?.decimals
      ),
    });
    depositTxHash.value = tx.hash;
    isDepositOngoing.value = true;
    await tx.wait();
  } else {
    const signer = transactor.getSigner();
    const tx = await signer.sendTransaction({
      to: L2PORTAL_ADDRESS,
      value: ethers.utils.parseUnits(
        Number(amount.value).toString(),
        currentChain.value.nativeCurrency?.decimals
      ),
    });
    withdrawTxHash.value = tx.hash;
    isWithdrawOngoing.value = true;
    await tx.wait();
  }
}

async function generateDepositProof(
  deposit: PendingDeposit
): Promise<MessageProof> {
  if (deposit.proofL1BlockNumber === undefined) {
    throw new Error("proofL1BlockNumber is undefined");
  }
  let rawProof = undefined;
  while (rawProof === undefined) {
    try {
      rawProof = await (l1Provider.value as JsonRpcProvider).send(
        "eth_getProof",
        [
          L1PORTAL_ADDRESS,
          [getStorageKey(deposit.depositHash)],
          ethers.utils.hexValue(deposit.proofL1BlockNumber),
        ]
      );
    } catch (e) {
      console.error(e);
    }
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  return {
    accountProof: rawProof.accountProof,
    storageProof: rawProof.storageProof[0].proof,
  };
}

async function generateWithdrawalProof(
  withdrawal: PendingWithdrawal
): Promise<MessageProof> {
  const rawProof = await (l2Provider.value as JsonRpcProvider).send(
    "eth_getProof",
    [
      L2PORTAL_ADDRESS,
      [getStorageKey(withdrawal.withdrawalHash)],
      ethers.utils.hexValue(withdrawal.l2BlockNumber),
    ]
  );
  return {
    accountProof: rawProof.accountProof,
    storageProof: rawProof.storageProof[0].proof,
  };
}

async function finalizeDeposit() {
  if (pendingDeposit.value === undefined) {
    return;
  }
  isDepositError.value = false;
  isDepositFinalizationRequested.value = true;
  await refreshBalance();
  const targetBalance = ethers.utils.parseEther(targetChainBalance.value.toString());
  if (DEPOSIT_BALANCE_THRESHOLD.gt(targetBalance)) {
    // request sequencer to help finalization
    try {
      const txHash = await requestFundDeposit(pendingDeposit.value);
      depositFinalizedTxHash.value = txHash;
    } catch (e) {
      console.error(e);
      isDepositFinalizationRequested.value = false;
      isDepositError.value = true;
    }
    return;
  }
  await switchNetworkAsync.value?.(targetChain.value.id);
  const l2Portal = IL2Portal__factory.connect(
    L2PORTAL_ADDRESS,
    transactor.getSigner(),
  );
  const l1Oracle = L1Oracle__factory.connect(
    L1ORACLE_ADDRESS,
    l1Provider.value,
  );
  try {
    const latestBlockNumber = await l1Oracle.blockNumber();
    pendingDeposit.value.proofL1BlockNumber = latestBlockNumber.toNumber();
    const proof = await generateDepositProof(pendingDeposit.value);
    const tx = await l2Portal.finalizeDepositTransaction(
      pendingDeposit.value.depositTx,
      proof.accountProof,
      proof.storageProof
    );
    depositFinalizedTxHash.value = tx.hash;
    await tx.wait();
  } catch (e) { 
    console.error(e);
    isDepositFinalizationRequested.value = false;
    isDepositError.value = true;
  }
  await switchNetworkAsync.value?.(currentChain.value.id);
}

async function finalizeWithdrawal() {
  if (
    pendingWithdrawal.value === undefined ||
    pendingWithdrawal.value.assertionID === undefined
  ) {
    return;
  }
  isWithdrawalError.value = false;
  isWithdrawalFinalizationRequested.value = true;
  await refreshBalance();
  await switchNetworkAsync.value?.(targetChain.value.id);
  const l1Portal = IL1Portal__factory.connect(
    L1PORTAL_ADDRESS,
    transactor.getSigner()
  );
  try {
    const proof = await generateWithdrawalProof(pendingWithdrawal.value);
    console.log(proof)
    console.log(pendingWithdrawal.value)
    let gas = await l1Portal.estimateGas.finalizeWithdrawalTransaction(
      pendingWithdrawal.value.withdrawalTx,
      pendingWithdrawal.value.assertionID,
      proof.accountProof,
      proof.storageProof,
      { gasLimit: 1000000 }, // avoid gas estimation error
    );
    console.log("gas", gas);
    gas = gas.add(150000); // extra gas to pass gas limit check in finalization
    const tx = await l1Portal.finalizeWithdrawalTransaction(
      pendingWithdrawal.value.withdrawalTx,
      pendingWithdrawal.value.assertionID,
      proof.accountProof,
      proof.storageProof,
      { gasLimit: gas },
    );
    withdrawFinalizedTxHash.value = tx.hash;
    await tx.wait();
  } catch (e) {
    console.error(e);
    isWithdrawalFinalizationRequested.value = false;
    isWithdrawalError.value = true;
    // pendingWithdrawal.value.assertionID = undefined;
    // pendingWithdrawal.value.proofL2BlockNumber = undefined;
    // isWithdrawalReadyToFinalize.value = false;
  }
  await switchNetworkAsync.value?.(currentChain.value.id);
}

async function closeDeposit() {
  pendingDeposit.value = undefined;
  depositTxHash.value = "";
  depositFinalizedTxHash.value = "";
  isDepositOngoing.value = false;
  isDepositReadyToFinalize.value = false;
  isDepositFinalizationRequested.value = false;
  isDepositFinished.value = false;
  isDepositError.value = false;
  await refreshBalance();
}

async function closeWithdraw() {
  pendingWithdrawal.value = undefined;
  withdrawTxHash.value = "";
  withdrawFinalizedTxHash.value = "";
  isWithdrawOngoing.value = false;
  isWithdrawalReadyToFinalize.value = false;
  isWithdrawalFinalizationRequested.value = false;
  isWithdrawFinished.value = false;
  isWithdrawalError.value = false;
  await refreshBalance();
}

onMounted(async () => {
  const inbox = ISequencerInbox__factory.connect(
    INBOX_ADDRESS,
    l1Provider.value
  );
  const rollup = IRollup__factory.connect(ROLLUP_ADDRESS, l1Provider.value);
  const l1Portal = IL1Portal__factory.connect(
    L1PORTAL_ADDRESS,
    l1Provider.value
  );
  const l1Oracle = L1Oracle__factory.connect(
    L1ORACLE_ADDRESS,
    l2Provider.value
  );
  const l2Portal = IL2Portal__factory.connect(
    L2PORTAL_ADDRESS,
    l2Provider.value
  );
  inbox.on(
    inbox.filters.TxBatchAppended(),
    async (batchNum, prevInboxSize, inboxSize, event) => {
      console.log("TxBatchAppended", batchNum.toString(), inboxSize.toString());
      if (pendingWithdrawal.value && !isWithdrawalReadyToFinalize.value) {
        if (pendingWithdrawal.value?.assertionID !== undefined) {
          // We already know which assertion this withdrawal is included in
          return;
        }
        // Get the last l2 block number of the current batch
        const tx = await event.getTransaction();
        const decoded = sequencerInboxInterface.decodeFunctionData(
          "appendTxBatch",
          tx.data
        );
        const contexts: BigNumber[] = decoded[0];
        const lastL2BlockNumber = contexts[contexts.length - 2].toNumber();
        console.log("L2BlockNumber", lastL2BlockNumber ,"<-> InboxSize", inboxSize.toString());
        // If it is larger than the pending withdrawal's l2 block number
        // The withdrawal is already sequenced on L1
        if (lastL2BlockNumber >= pendingWithdrawal.value.l2BlockNumber) {
          if (pendingWithdrawal.value?.inboxSize === undefined) {
            pendingWithdrawal.value.inboxSize = inboxSize;
          }
          inboxSizeToBlockNumberMap.set(
            inboxSize.toString(),
            lastL2BlockNumber
          );
        }
      }
    }
  );
  rollup.on(
    rollup.filters.AssertionCreated(),
    async (assertionID, asserter, vmHash, event) => {
      console.log("AssertionCreated", assertionID.toString());
      if (pendingWithdrawal.value && !isWithdrawalReadyToFinalize.value) {
        if (pendingWithdrawal.value?.inboxSize === undefined) {
          // We haven't seen the withdrawal sequenced on L1 yet
          return;
        }
        if (pendingWithdrawal.value?.assertionID !== undefined) {
          // We already know which assertion this withdrawal is included in
          return;
        }
        const assertion = await rollup.getAssertion(assertionID);
        console.log("Assertion ID", assertionID.toString(), "<-> InboxSize", assertion.inboxSize.toString());
        if (assertion.inboxSize.gte(pendingWithdrawal.value.inboxSize)) {
          // The assertion contains the withdrawal
          if (inboxSizeToBlockNumberMap.has(assertion.inboxSize.toString())) {
            // We already know the l2 block number of the assertion
            pendingWithdrawal.value.assertionID = assertionID;
            pendingWithdrawal.value.proofL2BlockNumber =
              inboxSizeToBlockNumberMap.get(assertion.inboxSize.toString());
            console.log(assertion.stateHash);
            // No need to keep the map
            inboxSizeToBlockNumberMap.clear();
          }
        }
      }
    }
  );
  rollup.on(rollup.filters.AssertionConfirmed(), async (assertionID, event) => {
    console.log("AssertionConfirmed", assertionID.toString());
    if (pendingWithdrawal.value && !isWithdrawalReadyToFinalize.value) {
      if (pendingWithdrawal.value?.assertionID === undefined) {
        return;
      }
      if (assertionID.gte(pendingWithdrawal.value.assertionID)) {
        // The assertion should be already finalized
        const assertion = await rollup.getAssertion(
          pendingWithdrawal.value.assertionID
        );
        if (assertion.inboxSize.eq(0)) {
          console.error("The assertion containing the withdrawal is rejected");
          console.error(
            "Assertion ID: ",
            pendingWithdrawal.value.assertionID.toString()
          );
          pendingWithdrawal.value.inboxSize = undefined;
          pendingWithdrawal.value.assertionID = undefined;
          pendingWithdrawal.value.proofL2BlockNumber = undefined;
          isWithdrawalError.value = true;
          return;
        }
        isWithdrawalReadyToFinalize.value = true;
      }
    }
  });
  l1Portal.on(
    l1Portal.filters.DepositInitiated(),
    (nonce, sender, target, value, gasLimit, data, depositHash, event) => {
      if (event.transactionHash === depositTxHash.value) {
        pendingDeposit.value = {
          l1BlockNumber: event.blockNumber,
          proofL1BlockNumber: undefined,
          depositHash: depositHash,
          depositTx: {
            nonce,
            sender,
            target,
            value,
            gasLimit,
            data,
          },
        };
      }
    }
  );
  l1Portal.on(
    l1Portal.filters.WithdrawalFinalized(),
    (_withdrawalHash, success, event) => {
      if (_withdrawalHash === pendingWithdrawal.value?.withdrawalHash) {
        isWithdrawFinished.value = true;
      }
    }
  );
  l1Oracle.on(
    l1Oracle.filters.L1OracleValuesUpdated(),
    (blockNumber, stateRoot, event) => {
      isDepositReadyToFinalize.value = false;
      if (pendingDeposit.value === undefined) {
        return;
      }
      if (blockNumber.gte(pendingDeposit.value.l1BlockNumber)) {
        isDepositReadyToFinalize.value = true;
        pendingDeposit.value.proofL1BlockNumber = blockNumber.toNumber();
      }
    }
  );
  l2Portal.on(
    l2Portal.filters.WithdrawalInitiated(),
    (nonce, sender, target, value, gasLimit, data, withdrawalHash, event) => {
      if (event.transactionHash === withdrawTxHash.value) {
        pendingWithdrawal.value = {
          l2BlockNumber: event.blockNumber,
          proofL2BlockNumber: undefined,
          inboxSize: undefined,
          assertionID: undefined,
          withdrawalHash: withdrawalHash,
          withdrawalTx: {
            nonce,
            sender,
            target,
            value,
            gasLimit,
            data,
          },
        };
      }
    }
  );
  l2Portal.on(
    l2Portal.filters.DepositFinalized(),
    (_depositHash, success, event) => {
      if (_depositHash === pendingDeposit.value?.depositHash) {
        isDepositFinished.value = true;
      }
    }
  );

  await refreshBalance();
});
</script>

<template>
  <v-card
    variant="tonal"
    rounded
  >
    <v-container>
      <v-row>
        <v-col col="12">
          <v-tabs
            v-model="mode"
            mandatory
            grow
          >
            <v-tab value="Deposit">
              Deposit
            </v-tab>
            <v-tab value="Withdraw">
              Withdraw
            </v-tab>
          </v-tabs>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-card rounded>
            <v-row justify="center">
              <v-col cols="11">
                <div class="pt-4 text-h4">
                  From <b>{{ currentChain.name }}</b>
                </div>
              </v-col>
            </v-row>
            <v-row justify="center">
              <v-col cols="11">
                <v-text-field
                  v-model="amount"
                  label="Amount"
                  variant="outlined"
                >
                  <template #append-inner>
                    <span>{{ currentChain.nativeCurrency?.symbol }}</span>
                  </template>
                </v-text-field>
                <div class="pb-4">
                  Balance:
                  {{
                    currentChainBalance
                  }}
                  {{ currentChain.nativeCurrency?.symbol }}
                </div>
              </v-col>
            </v-row>
          </v-card>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col
          cols="12"
          class="text-center"
        >
          <v-icon icon="mdi-arrow-down" />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-card rounded>
            <v-row justify="center">
              <v-col cols="11">
                <div class="pt-4 text-h4">
                  To <b>{{ targetChain.name }}</b>
                </div>
                <div class="pt-4">
                  You will receive: {{ amount }}
                  {{ targetChain.nativeCurrency?.symbol }}
                </div>
                <div class="pt-4 pb-4">
                  Balance:
                  {{
                    targetChainBalance
                  }}
                  {{ targetChain.nativeCurrency?.symbol }}
                </div>
              </v-col>
            </v-row>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col col="12">
          <v-btn
            v-if="!isConnected"
            block
            @click="async () => { await connect(); await refreshBalance();}"
          >
            Connect Wallet
          </v-btn>
          <v-btn
            v-else-if="chain?.id !== currentChain.id"
            block
            @click="
              async () => {
                await switchNetworkAsync?.(currentChain.id);
                await connect();
                await refreshBalance();
              }
            "
          >
            Switch To {{ currentChain.name }}
          </v-btn>
          <template v-else>
            <v-btn
              block
              :disabled="isInsufficientBalance || isZeroAmount || isNaNAmount"
              :title="bridgeBtnText"
              @click="bridge"
            >
              {{ mode }}
            </v-btn>
            <v-btn
              block
              class="mt-4"
              @click="refreshBalance"
            >
              Refresh Balance
            </v-btn>
          </template>
        </v-col>
      </v-row>
    </v-container>
    <v-dialog
      v-model="isDepositOngoing"
      width="auto"
      persistent
    >
      <v-card>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-alert
                v-if="isDepositError"
                color="error"
                icon="$error"
                title="Deposit Finalization Failed"
                text="Please try again. If the problem persists, please record the deposit initiation transaction hash for recovery."
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <div class="text-h4 pb-2">
                Deposit Details
              </div>
              <div class="text-h6">
                Deposit initiation:
              </div>
              <div>
                <a
                  v-if="depositTxHash !== ''"
                  :href="`${currentChain.blockExplorers?.default.url}/tx/${depositTxHash}`"
                  target="_blank"
                >{{ depositTxHash }}</a>
                <span v-else>Pending...</span>
              </div>
              <div class="text-h6">
                Deposit finalization:
              </div>
              <div>
                <a
                  v-if="depositFinalizedTxHash !== ''"
                  :href="`${targetChain.blockExplorers?.default.url}/tx/${depositFinalizedTxHash}`"
                  target="_blank"
                >{{ depositFinalizedTxHash }}</a>
                <span v-else>Pending...</span>
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn
                :disabled="
                  !isDepositReadyToFinalize || isDepositFinalizationRequested
                "
                :color="isDepositFinished ? 'success' : ''"
                @click="finalizeDeposit"
              >
                Finalize Deposit
              </v-btn>
              <v-btn
                v-if="isDepositFinished"
                @click="closeDeposit"
              >
                Finish
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-dialog>
    <v-dialog
      v-model="isWithdrawOngoing"
      width="auto"
      persistent
    >
      <v-card>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-alert
                v-if="isWithdrawalError"
                color="error"
                icon="$error"
                title="Withdrawal Finalization Failed"
                text="Please try again. If the problem persists, please record the withdrawal initiation transaction hash for recovery."
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <div class="text-h4 pb-2">
                Withdrawal Details
              </div>
              <div class="text-h6">
                Withdrawal initiation:
              </div>
              <div>
                <a
                  v-if="withdrawTxHash !== ''"
                  :href="`${currentChain.blockExplorers?.default.url}/tx/${withdrawTxHash}`"
                  target="_blank"
                >{{ withdrawTxHash }}</a>
                <span v-else>Pending...</span>
              </div>
              <div class="text-h6">
                Withdrawal finalization:
              </div>
              <div>
                <a
                  v-if="withdrawFinalizedTxHash !== ''"
                  :href="`${targetChain.blockExplorers?.default.url}/tx/${withdrawFinalizedTxHash}`"
                  target="_blank"
                >{{ withdrawFinalizedTxHash }}</a>
                <span v-else>Pending...</span>
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn
                :disabled="
                  !isWithdrawalReadyToFinalize ||
                    isWithdrawalFinalizationRequested
                "
                :color="isWithdrawFinished ? 'success' : ''"
                @click="finalizeWithdrawal"
              >
                Finalize Withdrawal
              </v-btn>
              <v-btn
                v-if="isWithdrawFinished"
                @click="closeWithdraw"
              >
                Finish
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-dialog>
  </v-card>
</template>
