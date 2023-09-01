import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-create-assertion.styles';
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

function TxCreateAssertion({ wallet, withdrawData,pendingWithdraw,setPendingWithdraw,switchChain,onGoBack, onGoToNextStep }: TxPendingProps) {
  const classes = useTxPendingStyles();
  const l1Provider = new ethers.providers.StaticJsonRpcProvider(CHIADO_RPC_URL);
  const inboxSizeToBlockNumberMap = new Map<string, number>();
  const rollup = IRollup__factory.connect(ROLLUP_ADDRESS, l1Provider);

  useEffect(() => {

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
            // We already know the l2 block number of the assertion
            pendingWithdraw.data.assertionID = assertionID;
            setPendingWithdraw(pendingWithdraw);
            console.log(assertion.stateHash);
            onGoToNextStep();
            // No need to keep the map
            inboxSizeToBlockNumberMap.clear();
        }
      }
    }
  );
  },[pendingWithdraw,rollup, withdrawData,rollup.filters.AssertionCreated()]
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
      Waiting for Assertion Creation
    </div>
  );
}

export default TxCreateAssertion;
