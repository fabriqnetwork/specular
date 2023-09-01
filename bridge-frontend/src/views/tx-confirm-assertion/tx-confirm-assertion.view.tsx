import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-confirm-assertion.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import { ethers } from 'ethers'
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

function TxConfirmAssertion({ wallet, withdrawData,pendingWithdraw,setPendingWithdraw,switchChain,onGoBack, onGoToNextStep }: TxPendingProps) {
  const classes = useTxPendingStyles();
  const l1Provider = new ethers.providers.StaticJsonRpcProvider(CHIADO_RPC_URL);
  const rollup = IRollup__factory.connect(ROLLUP_ADDRESS, l1Provider);

  useEffect(() => {
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
          onGoBack();
          return;
        }
        setPendingWithdraw(pendingWithdraw);
        onGoToNextStep();
      }
    }
  });
  },[pendingWithdraw,rollup]
  )
  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`Specular Bridge`} />
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
      Waiting for Assertion Confirmation
    </div>
  );
}

export default TxConfirmAssertion;
