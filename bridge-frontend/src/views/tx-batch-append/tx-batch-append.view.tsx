import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-batch-append.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import { ethers, BigNumber } from 'ethers'
import {
  CHIADO_RPC_URL,
  INBOX_ADDRESS,
  ROLLUP_ADDRESS,
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
  onGoToNextStep: () => void;
  onGoBack: () => void;
  switchChain: (args1:any) => void;

}

function TxBatchAppend({ wallet, withdrawData,pendingWithdraw,setPendingWithdraw,switchChain,onGoBack, onGoToNextStep }: TxPendingProps) {
  const classes = useTxPendingStyles();
  const l1Provider = new ethers.providers.StaticJsonRpcProvider(CHIADO_RPC_URL);
  const sequencerInboxInterface = ISequencerInbox__factory.createInterface();
  const inbox = ISequencerInbox__factory.connect(
    INBOX_ADDRESS,
    l1Provider
  );

  // useEffect(() => {

    inbox.on(
      inbox.filters.TxBatchAppended(),
      async (batchNum, prevInboxSize, inboxSize, event) => {
        console.log("TxBatchAppended", batchNum.toString(), inboxSize.toString());
        if (pendingWithdraw ) {
          // Get the last l2 block number of the current batch
          const tx = await event.getTransaction();
          const decoded = sequencerInboxInterface.decodeFunctionData(
            "appendTxBatch",
            tx.data
          );
          const contexts: BigNumber[] = decoded[0];
          const firstL2BlockNumber = decoded[2];
          const lastL2BlockNumber = contexts.length / 2 + firstL2BlockNumber.toNumber() - 1;
          console.log("L2BlockNumber", lastL2BlockNumber ,"<-> InboxSize", inboxSize.toString());
          // If it is larger than the pending withdrawal's l2 block number
          // The withdrawal is already sequenced on L1
          if (lastL2BlockNumber >= pendingWithdraw.data.l2BlockNumber) {
            if (pendingWithdraw.data?.inboxSize === undefined) {
              pendingWithdraw.data.inboxSize = inboxSize;
              pendingWithdraw.data.proofL2BlockNumber = lastL2BlockNumber;
              setPendingWithdraw(pendingWithdraw);
            }
            onGoToNextStep();
          }
        }
      }
    );

  // },[withdrawData,pendingWithdraw,inbox,inbox.filters.TxBatchAppended()])

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
      Waiting for Transaction to be processed
    </div>
  );
}

export default TxBatchAppend;
