import useTxConfirmStyles from './tx-confirm.styles';
import Header from '../shared/header/header.view';
import Spinner from '../shared/spinner/spinner.view';
import { useEffect } from 'react';

interface TxConfirmProps {
  wallet: {
    address: string;
  };
  depositData: {
    status: string;
  };
  onGoBack: () => void;
  onGoToPendingStep: () => void;
}

function TxConfirm({ wallet, depositData, onGoBack, onGoToPendingStep }: TxConfirmProps) {
  const classes = useTxConfirmStyles();

  useEffect(() => {
    if (depositData.status === 'failed') {
      onGoBack();
    }
    if (depositData.status === 'pending') {
      onGoToPendingStep();
    }
  }, [depositData, onGoBack, onGoToPendingStep]);

  return (
    <div className={classes.txConfirm}>
      <Header address={wallet.address} title={`xDAI → ETH`} />
      <div className={classes.spinnerWrapper}>
        <Spinner />
      </div>
      <p className={classes.title}>Confirm the transaction in your wallet</p>
    </div>
  );
}

export default TxConfirm;
