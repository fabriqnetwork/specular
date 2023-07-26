import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-pending-finalize-withdraw.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import { ethers } from 'ethers'
import * as React from 'react';
import {
  L1PORTAL_ADDRESS
} from "../../constants";
import {
  IL2Portal__factory,
} from "../../typechain-types";
import type { PendingWithdrawal } from "../../types";

interface TxPendingProps {
 wallet: {
    address: string;
    chainId: number;
    provider: any;
  };
  depositData: {
    status: string;
    data?: {
      hash: string;
    };
  };
  setPendingWithdraw: (arg0:PendingWithdrawal) => void;
  onGoToFinalizeStep: () => void;

}

function TxPendingFinalizeWithdraw({ wallet, depositData, setPendingWithdraw, onGoToFinalizeStep }: TxPendingProps) {
  const classes = useTxPendingStyles();

  const l2Provider =  wallet.provider
  const l2Portal = IL2Portal__factory.connect(
    L1PORTAL_ADDRESS,
    l2Provider
  );

  useEffect(() => {
    l2Portal.on(
      l2Portal.filters.WithdrawalInitiated(),
      (nonce:ethers.BigNumber, sender:string, target:string, value:ethers.BigNumber, gasLimit:ethers.BigNumber, data:string, withdrawalHash:string, event:any) => {
        if (event.transactionHash === depositData.data?.hash) {
          const newPendingWithdrawal: PendingWithdrawal = {
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
          }
          setPendingWithdraw(newPendingWithdrawal);
          onGoToFinalizeStep();
      }
    }
  )
  }, [l2Portal, depositData, setPendingWithdraw, onGoToFinalizeStep]);

  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`xDAI → ETH`} />
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[wallet.chainId].blockExplorerUrl}/tx/${depositData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction status here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
    </div>
  );
}

export default TxPendingFinalizeWithdraw;
