import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-pending.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import * as React from 'react';

interface TxPendingProps {
  wallet: {
    address: string;
    chainId: number;
    provider: any;
  };
  transactionData: {
    status: string;
    data?: {
      hash: string;
    };
  };
  onGoBack: () => void;
  onGoToFinalizeStep: () => void;
}

function TxPending({ wallet, transactionData, onGoBack, onGoToFinalizeStep }: TxPendingProps) {
  const classes = useTxPendingStyles();

  useEffect(() => {
    if (transactionData.status === 'failed') {
      onGoBack();
    }
    if (transactionData.status === 'successful') {
      onGoToFinalizeStep();
    }
  }, [transactionData, onGoBack, onGoToFinalizeStep]);

  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`xDAI → ETH`} />
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[wallet.chainId].blockExplorerUrl}/tx/${transactionData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction status here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      Proceeding to Finalize Transaction
    </div>
  );
}

export default TxPending;
