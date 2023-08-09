import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-pending-finalize-withdraw.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import { ethers } from 'ethers'
import {
  CHIADO_NETWORK_ID,
  CHIADO_RPC_URL,
  INBOX_ADDRESS,
  ROLLUP_ADDRESS,
  SPECULAR_RPC_URL,
  SPECULAR_NETWORK_ID
} from "../../constants";
import {
  ISequencerInbox__factory,
  IRollup__factory
} from "../../typechain-types";
import type { PendingWithdrawal } from "../../types";

interface PendingWithdrawlData {
  status: string;
  data: PendingWithdrawal;
}
interface TxPendingProps {
 wallet: {
    address: string;
    chainId: number;
    provider: any;
  };
  withdrawData: {
    status: string;
    data?: {
      hash: string;
    };
  };
  pendingWithdraw: PendingWithdrawlData;
  setPendingWithdraw: (args1:any) => void;
  onGoToFinalizeStep: () => void;
  switchChain: (args1:any) => void;

}

function TxPendingFinalizeWithdraw({ wallet, withdrawData,pendingWithdraw,setPendingWithdraw,switchChain, onGoToFinalizeStep }: TxPendingProps) {
  const classes = useTxPendingStyles();
  const l1Provider = new ethers.providers.StaticJsonRpcProvider(CHIADO_RPC_URL);
  const l2Provider = new ethers.providers.StaticJsonRpcProvider(SPECULAR_RPC_URL);
  const sequencerInboxInterface = ISequencerInbox__factory.createInterface();
  const inboxSizeToBlockNumberMap = new Map<string, number>();

  useEffect(() => {
    const inbox = ISequencerInbox__factory.connect(
      INBOX_ADDRESS,
      l1Provider
    );
    const rollup = IRollup__factory.connect(ROLLUP_ADDRESS, l1Provider);

    inbox.on(
      inbox.filters.TxBatchAppended(),
      async (batchNum, prevInboxSize, inboxSize, event) => {
        console.log("TxBatchAppended", batchNum.toString(), inboxSize.toString());
        if (pendingWithdraw ) {
          if (pendingWithdraw.data?.assertionID !== undefined) {
            // We already know which assertion this withdrawal is included in
            return;
          }
          // Get the last l2 block number of the current batch
          const tx = await event.getTransaction();
          const decoded = sequencerInboxInterface.decodeFunctionData(
            "appendTxBatch",
            tx.data
          );
          const contexts: ethers.BigNumber[] = decoded[0];
          const lastL2BlockNumber = contexts[contexts.length - 2].toNumber();
          console.log("L2BlockNumber", lastL2BlockNumber ,"<-> InboxSize", inboxSize.toString());
          // If it is larger than the pending withdrawal's l2 block number
          // The withdrawal is already sequenced on L1
          if (lastL2BlockNumber >= pendingWithdraw.data.l2BlockNumber) {
            if (pendingWithdraw.data?.inboxSize === undefined) {
              pendingWithdraw.data.inboxSize = inboxSize;
              setPendingWithdraw(pendingWithdraw);
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
      if (pendingWithdraw ) {
        if (pendingWithdraw.data?.inboxSize === undefined) {
          // We haven't seen the withdrawal sequenced on L1 yet
          return;
        }
        if (pendingWithdraw.data?.assertionID !== undefined) {
          // We already know which assertion this withdrawal is included in
          return;
        }
        const assertion = await rollup.getAssertion(assertionID);
        console.log("Assertion ID", assertionID.toString(), "<-> InboxSize", assertion.inboxSize.toString());
        if (assertion.inboxSize.gte(pendingWithdraw.data.inboxSize)) {
          // The assertion contains the withdrawal
          if (inboxSizeToBlockNumberMap.has(assertion.inboxSize.toString())) {
            // We already know the l2 block number of the assertion
            pendingWithdraw.data.assertionID = assertionID;
            pendingWithdraw.data.proofL2BlockNumber = inboxSizeToBlockNumberMap.get(assertion.inboxSize.toString());
            setPendingWithdraw(pendingWithdraw);
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
    if (pendingWithdraw) {
      if (pendingWithdraw.data.assertionID === undefined) {
        return;
      }
      console.log("pendingWithdrawal assertionID ", pendingWithdraw.data.assertionID.toString());
      if (assertionID.gte(pendingWithdraw.data.assertionID)) {
        // The assertion should be already finalized
        const assertion = await rollup.getAssertion(
          pendingWithdraw.data.assertionID
        );
        if (assertion.inboxSize.eq(0)) {
          console.error("The assertion containing the withdrawal is rejected");
          console.error("Assertion ID: ", pendingWithdraw.data.assertionID.toString()
          );
          pendingWithdraw.data.inboxSize = undefined;
          pendingWithdraw.data.assertionID = undefined;
          pendingWithdraw.data.proofL2BlockNumber = undefined;
          // isWithdrawalError.value = true;
          return;
        }
        setPendingWithdraw(pendingWithdraw);
        onGoToFinalizeStep();
      }
    }
  });
  },[pendingWithdraw]
  )
  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`xDAI → ETH`} />
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[SPECULAR_NETWORK_ID].blockExplorerUrl}/tx/${withdrawData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Transaction on {NETWORKS[SPECULAR_NETWORK_ID].name} is successful. Check transaction details here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      Waiting for confirmation
    </div>
  );
}

export default TxPendingFinalizeWithdraw;
