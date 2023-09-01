import useTxConfirmStyles from './tx-confirm.styles';
import Header from '../shared/header/header.view';
import * as React from 'react';
import Spinner from '../shared/spinner/spinner.view';
import { useEffect } from 'react';

interface TxConfirmProps {
  wallet: {
    address: string;
    chainId: number;
    provider: any;
  };
  transactionData: {
    status: string;
  };
  onGoBack: () => void;
  onGoToNextStep: () => void;
}

function TxConfirm({ wallet, transactionData, onGoBack, onGoToNextStep }: TxConfirmProps) {
  const classes = useTxConfirmStyles();

  useEffect(() => {

    if (transactionData.status === 'failed') {
      onGoBack();
    }
    if (transactionData.status === 'pending') {
      onGoToNextStep();
    }
  }, [transactionData, onGoBack, onGoToNextStep]);

  return (
    <div className={classes.txConfirm}>
      <Header address={wallet.address} title={`Specular Bridge`} />
      <div className={classes.spinnerWrapper}>
        <Spinner />
      </div>
      <p className={classes.title}>Confirm the transaction in your wallet</p>
    </div>
  );
}

export default TxConfirm;
