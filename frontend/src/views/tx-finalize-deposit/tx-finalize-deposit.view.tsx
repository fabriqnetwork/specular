import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxFinalizeDepositStyles from './tx-finalize-deposit.styles';
import { ReactComponent as LinkIcon } from '../../images/link-icon.svg';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';
import {
  CHIADO_NETWORK_ID,
  SPECULAR_NETWORK_ID
} from "../../constants";

interface TxFinalizeDepositProps {
  wallet: {
    address: string;
    chainId: number;
  };
  depositData: {
    status: string;
    data?: {
      hash: string;
    };
  };
  finalizeDepositData: {
    status: string;
    data?: string;
  };
  onGoBack: () => void;
  onGoToOverviewStep: () => void;
}

function TxFinalizeDeposit({ wallet, depositData, finalizeDepositData, onGoBack, onGoToOverviewStep }: TxFinalizeDepositProps) {
  const classes = useTxFinalizeDepositStyles();

  useEffect(() => {
    if (depositData.status === 'failed') {
      onGoBack();
    }
    if (depositData.status === 'successful') {
      onGoToOverviewStep();
    }
  }, [depositData, onGoBack, onGoToOverviewStep]);

  return (
    <div className={classes.txOverview}>
      <Header address={wallet.address} title={`xDAI → ETH`} />
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[CHIADO_NETWORK_ID].blockExplorerUrl}/tx/${depositData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check transaction details for deposit here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[SPECULAR_NETWORK_ID].blockExplorerUrl}/tx/${finalizeDepositData?.data}`}
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

export default TxFinalizeDeposit;
