import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-pending-withdraw.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import {
  L2PORTAL_ADDRESS,
  SPECULAR_NETWORK_ID
} from "../../constants";
import type {PendingWithdrawal } from "../../types";
import {
  IL2Portal__factory,
} from "../../typechain-types";
import { ethers } from 'ethers';

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
  l2Provider : any;
  pendingWithdraw : PendingWithdrawlData;
  setPendingWithdraw: (args1:any) => void;
  onGoBack: () => void;
  onGoToFinalizeStep: () => void;
}

function TxPendingWithdraw({ wallet, withdrawData, l2Provider, pendingWithdraw, setPendingWithdraw, onGoBack, onGoToFinalizeStep }: TxPendingProps) {
  const classes = useTxPendingStyles();


  useEffect(() => {
    if (withdrawData.status === 'failed') {
      onGoBack();
    }
    if (pendingWithdraw.status==='initiated') {
      onGoToFinalizeStep();
    }
  }, [withdrawData, pendingWithdraw, onGoBack, onGoToFinalizeStep]);

    useEffect(() => {
    const l2Portal = IL2Portal__factory.connect(
      L2PORTAL_ADDRESS,
      l2Provider
    );
    const version = 0;
    l2Portal.on(
      l2Portal.filters.WithdrawalInitiated(),
      (nonce, sender, target, value, gasLimit, data, withdrawalHash, event) => {
        if (event.transactionHash === withdrawData.data?.hash) {
          const newPendingWithdrawal: PendingWithdrawal ={
            l2BlockNumber: event.blockNumber,
            proofL2BlockNumber: undefined,
            inboxSize: undefined,
            assertionID: undefined,
            withdrawalHash: withdrawalHash,
            withdrawalTx: {
              version,
              nonce,
              sender,
              target,
              value,
              gasLimit,
              data,
            },
          }
          console.log("Main Correct L1 Portal transactionHash is "+event.transactionHash)
          setPendingWithdraw({ status: 'initiated', data: newPendingWithdrawal});
        }
      }
    );
    },[withdrawData]
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
          href={`${NETWORKS[wallet.chainId].blockExplorerUrl}/tx/${withdrawData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check {NETWORKS[SPECULAR_NETWORK_ID].name}'s transaction status here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      Proceeding to Finalize Transaction
    </div>
  );
}

export default TxPendingWithdraw;
