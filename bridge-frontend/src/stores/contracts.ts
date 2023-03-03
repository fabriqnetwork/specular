import {
  INBOX_ADDRESS,
  L1ORACLE_ADDRESS,
  L1PORTAL_ADDRESS,
  L2PORTAL_ADDRESS,
  ROLLUP_ADDRESS,
} from "@/client/constants";
import type { Provider } from "@ethersproject/providers";
import { defineStore } from "pinia";
import {
  IL1Portal__factory,
  IL2Portal__factory,
  IRollup__factory,
  ISequencerInbox__factory,
  L1Oracle__factory,
  type IL1Portal,
  type IL2Portal,
  type IRollup,
  type ISequencerInbox,
  type L1Oracle,
} from "../../../contracts/typechain-types";

export const useContractsStore = defineStore("contracts", {
  state: () => ({
    inbox: null as ISequencerInbox | null,
    rollup: null as IRollup | null,
    l1Portal: null as IL1Portal | null,
    l2Portal: null as IL2Portal | null,
    l1Oracle: null as L1Oracle | null,
  }),
  actions: {
    init(l1Provider: Provider, l2Provider: Provider) {
      this.inbox = ISequencerInbox__factory.connect(INBOX_ADDRESS, l1Provider);
      this.rollup = IRollup__factory.connect(ROLLUP_ADDRESS, l1Provider);
      this.l1Portal = IL1Portal__factory.connect(L1PORTAL_ADDRESS, l1Provider);
      this.l2Portal = IL2Portal__factory.connect(L2PORTAL_ADDRESS, l2Provider);
      this.l1Oracle = L1Oracle__factory.connect(L1ORACLE_ADDRESS, l2Provider);
    },
  },
});
