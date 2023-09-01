import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxFinalizeWithdrawStyles from './tx-finalize-withdraw.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import {
  CHIADO_NETWORK_ID,
  SPECULAR_NETWORK_ID
} from "../../constants";

interface TxFinalizeWithdrawProps {
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
  finalizeWithdrawData: {
    status: string;
    data?: string;
  };
  onGoBack: () => void;
  onGoToOverviewStep: () => void;
}

function TxFinalizeWithdraw({ wallet, withdrawData, finalizeWithdrawData, onGoBack, onGoToOverviewStep }: TxFinalizeWithdrawProps) {
  const classes = useTxFinalizeWithdrawStyles();

  useEffect(() => {
    if (finalizeWithdrawData.status === 'failed') {
      onGoBack();
    }
    if (finalizeWithdrawData.status === 'successful') {
      onGoToOverviewStep();
    }
  }, [finalizeWithdrawData, onGoBack, onGoToOverviewStep]);

  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`Specular Bridge`} />
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[SPECULAR_NETWORK_ID].blockExplorerUrl}/tx/${withdrawData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction details for withdraw here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[CHIADO_NETWORK_ID].blockExplorerUrl}/tx/${finalizeWithdrawData?.data}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction finalization status here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
    </div>
  );
}

export default TxFinalizeWithdraw;
