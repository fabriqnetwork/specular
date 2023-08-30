import useTxConfirmChainStyles from './tx-confirm-chain.styles';
import Header from '../shared/header/header.view';
import * as React from 'react';
import Spinner from '../shared/spinner/spinner.view';
import { useEffect } from 'react';

interface TxConfirmChainProps {
  wallet: {
    address: string;
    chainId: number;
    provider: any;
  };
  networkId:string;
  onGoBack: () => void;
  onGoToNextStep: () => void;
}

function TxConfirmChain({ wallet, networkId, onGoBack, onGoToNextStep }: TxConfirmChainProps) {
  console.log("Network Id is",networkId);
  const classes = useTxConfirmChainStyles();

  useEffect(() => {

    console.log("Chain Id is",wallet.chainId.toString());

    if (wallet.chainId.toString() === networkId) {
      onGoToNextStep();
    }
  }, [wallet.chainId, networkId, onGoBack, onGoToNextStep]);

  return (
    <div className={classes.txConfirmChain}>
      <Header address={wallet.address} title={`xDAI → ETH`} />
      <div className={classes.spinnerWrapper}>
        <Spinner />
      </div>
      <p className={classes.title}>Confirm the transaction in your wallet</p>
    </div>
  );
}

export default TxConfirmChain;
