import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-pending-deposit.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import {
  L1PORTAL_ADDRESS,
  CHIADO_NETWORK_ID
} from "../../constants";
import type { PendingDeposit } from "../../types";
import {
  IL1Portal__factory,
} from "../../typechain-types";
import { ethers } from 'ethers';

interface PendingData {
  status: string;
  data: PendingDeposit;
}

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
  l1Provider : any;
  pendingDeposit : PendingData;
  setPendingDeposit: (args1:any) => void;
  onGoBack: () => void;
  onGoToFinalizeStep: () => void;
}

function TxPendingDeposit({ wallet, depositData, l1Provider, pendingDeposit, setPendingDeposit, onGoBack, onGoToFinalizeStep }: TxPendingProps) {
  const classes = useTxPendingStyles();


  useEffect(() => {
    if (depositData.status === 'failed') {
      onGoBack();
    }
    if (depositData.status === 'successful' && pendingDeposit.status==='initiated') {
      onGoToFinalizeStep();
    }
  }, [depositData, pendingDeposit, onGoBack, onGoToFinalizeStep]);

    useEffect(() => {
    const l1Portal = IL1Portal__factory.connect(
      L1PORTAL_ADDRESS,
      l1Provider
    );
    l1Portal.on(
      l1Portal.filters.DepositInitiated(),
      (nonce:ethers.BigNumber, sender:string, target:string, value:ethers.BigNumber, gasLimit:ethers.BigNumber, data:string, depositHash:string, event:any) => {
        console.log("Main L1 Portal transactionHash is "+event.transactionHash+"Deposit data hash is "+ depositData.status+" Deposit Data is "+depositData.data+" Deposit Data hash is "+depositData.data?.hash)
        if (event.transactionHash === depositData.data?.hash) {
          const newPendingDeposit: PendingDeposit = {
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
          }
          console.log("Main Correct L1 Portal transactionHash is "+event.transactionHash)
          setPendingDeposit({ status: 'initiated', data: newPendingDeposit});
        }
      }
    );
    },[depositData]
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
          href={`${NETWORKS[wallet.chainId].blockExplorerUrl}/tx/${depositData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check {NETWORKS[CHIADO_NETWORK_ID].name}'s transaction status here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      Proceeding to Finalize Transaction
    </div>
  );
}

export default TxPendingDeposit;
