<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
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
import { L1PORTAL_ADDRESS, L2PORTAL_ADDRESS } from "@/client/constants";
import { transactor } from "@/client/providers";
import {
  IL1Portal__factory,
  IL2Portal__factory,
} from "../../../contracts/typechain-types";

const { address, isConnected } = useAccount();
const { data: ensName } = useEnsName({
  address,
});
const { connect } = useConnect({
  connector: new InjectedConnector(),
});
const { chains, error, isLoading, pendingChainId, switchNetwork } =
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
const currentChainBalance = ref(0);
const targetChainBalance = ref(0);

async function refreshL1Balance() {
  if (address.value) {
    const l1Balance = await l1Provider.value.getBalance(address.value);
    if (mode.value === "Deposit") {
      currentChainBalance.value = Number(
        ethers.utils.formatUnits(
          l1Balance,
          currentChain.value.nativeCurrency?.decimals
        )
      );
    } else {
      targetChainBalance.value = Number(
        ethers.utils.formatUnits(
          l1Balance,
          targetChain.value.nativeCurrency?.decimals
        )
      );
    }
  }
}

async function refreshL2Balance() {
  if (address.value) {
    const l2Balance = await l2Provider.value.getBalance(address.value);
    if (mode.value === "Deposit") {
      targetChainBalance.value = Number(
        ethers.utils.formatUnits(
          l2Balance,
          targetChain.value.nativeCurrency?.decimals
        )
      );
    } else {
      currentChainBalance.value = Number(
        ethers.utils.formatUnits(
          l2Balance,
          currentChain.value.nativeCurrency?.decimals
        )
      );
    }
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

onMounted(async () => {
  const l1Portal = IL1Portal__factory.connect(
    L1PORTAL_ADDRESS,
    l1Provider.value
  );
  const l2Portal = IL2Portal__factory.connect(
    L2PORTAL_ADDRESS,
    l2Provider.value
  );
  l1Portal.on(
    l1Portal.filters.DepositInitiated(),
    (nonce, sender, target, value, gasLimit, data, _depositHash, event) => {
      if (event.transactionHash === depositTxHash.value) {
        depositHash.value = _depositHash;
      }
    }
  );
  l1Portal.on(
    l1Portal.filters.WithdrawalFinalized(),
    (_withdrawalHash, success, event) => {
      if (_withdrawalHash === withdrawHash.value) {
        isWithdrawFinished.value = true;
        withdrawFinalizedTxHash.value = event.transactionHash;
      }
    }
  );
  l2Portal.on(
    l2Portal.filters.WithdrawalInitiated(),
    (nonce, sender, target, value, gasLimit, data, _withdrawal, event) => {
      if (event.transactionHash === withdrawTxHash.value) {
        withdrawHash.value = _withdrawal;
      }
    }
  );
  l2Portal.on(
    l2Portal.filters.DepositFinalized(),
    (_depositHash, success, event) => {
      if (_depositHash === depositHash.value) {
        isDepositFinished.value = true;
        depositFinalizedTxHash.value = event.transactionHash;
      }
    }
  );

  await refreshBalance();
});

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
    const signer = await transactor.getSigner();
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
    const signer = await transactor.getSigner();
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
const depositFinalizedTxHash = ref("");
const withdrawFinalizedTxHash = ref("");
const depositHash = ref("");
const withdrawHash = ref("");
const isDepositOngoing = ref(false);
const isWithdrawOngoing = ref(false);
const isDepositFinished = ref(false);
const isWithdrawFinished = ref(false);

function closeDeposit() {
  isDepositOngoing.value = false;
  isDepositFinished.value = false;
}

function closeWithdraw() {
  isWithdrawOngoing.value = false;
  isWithdrawFinished.value = false;
}
</script>

<template>
  <v-card variant="tonal" rounded>
    <v-container>
      <v-row>
        <v-col col="12">
          <v-tabs v-model="mode" mandatory grow>
            <v-tab value="Deposit">Deposit</v-tab>
            <v-tab value="Withdraw">Withdraw</v-tab>
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
                  label="Amount"
                  variant="outlined"
                  v-model="amount"
                >
                  <template v-slot:append-inner>
                    <span>{{ currentChain.nativeCurrency?.symbol }}</span>
                  </template>
                </v-text-field>
                <div class="pb-4">
                  Balance: {{ currentChainBalance }}
                  {{ currentChain.nativeCurrency?.symbol }}
                </div>
              </v-col>
            </v-row>
          </v-card>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="12" class="text-center">
          <v-icon icon="mdi-arrow-down"></v-icon>
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
                  Balance: {{ targetChainBalance }}
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
            block
            v-if="chain?.id !== currentChain.id"
            @click="() => switchNetwork?.(currentChain.id)"
            >Switch To {{ currentChain.name }}</v-btn
          >
          <v-btn block v-else-if="!isConnected" @click="() => connect()"
            >Connect Wallet</v-btn
          >
          <template v-else>
            <v-btn
              @click="bridge"
              block
              :disabled="isInsufficientBalance || isZeroAmount || isNaNAmount"
              :title="bridgeBtnText"
              >{{ mode }}</v-btn
            >
            <v-btn @click="refreshBalance" block class="mt-4"
              >Refresh Balance</v-btn
            >
          </template>
        </v-col>
      </v-row>
    </v-container>
    <v-dialog width="auto" persistent v-model="isDepositOngoing">
      <v-card>
        <v-container>
          <v-row>
            <v-col cols="12">
              <div class="text-h4 pb-2">Deposit Details</div>
              <div class="text-h6">Deposit initiation:</div>
              <div>
                <a
                  v-if="depositTxHash !== ''"
                  :href="`${currentChain.blockExplorers?.default.url}/tx/${depositTxHash}`"
                  target="_blank"
                  >{{ depositTxHash }}</a
                >
                <span v-else>Pending...</span>
              </div>
              <div class="text-h6">Deposit finalization:</div>
              <div>
                <a
                  v-if="depositFinalizedTxHash !== ''"
                  :href="`${targetChain.blockExplorers?.default.url}/tx/${depositFinalizedTxHash}`"
                  target="_blank"
                  >{{ depositFinalizedTxHash }}</a
                >
                <span v-else>Pending...</span>
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn @click="closeDeposit" v-if="isDepositFinished"></v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-dialog>
    <v-dialog width="auto" persistent v-model="isWithdrawOngoing">
      <v-card>
        <v-container>
          <v-row>
            <v-col cols="12">
              <div class="text-h4 pb-2">Withdrawal Details</div>
              <div class="text-h6">Withdrawal initiation:</div>
              <div>
                <a
                  v-if="withdrawTxHash !== ''"
                  :href="`${currentChain.blockExplorers?.default.url}/tx/${withdrawTxHash}`"
                  target="_blank"
                  >{{ withdrawHash }}</a
                >
                <span v-else>Pending...</span>
              </div>
              <div class="text-h6">Withdrawal finalization:</div>
              <div>
                <a
                  v-if="withdrawFinalizedTxHash !== ''"
                  :href="`${targetChain.blockExplorers?.default.url}/tx/${withdrawFinalizedTxHash}`"
                  target="_blank"
                  >{{ withdrawFinalizedTxHash }}</a
                >
                <span v-else>Pending...</span>
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn @click="closeWithdraw" v-if="isWithdrawFinished"></v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-dialog>
  </v-card>
</template>
