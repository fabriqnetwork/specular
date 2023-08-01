import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxFinalizeDepositStyles from './tx-finalize-deposit.styles';
import LinkIcon from '@mui/icons-material/OpenInNew';
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
    provider: any;
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
    // if (finalizeDepositData.status === 'failed') {
    //   onGoBack();
    // }
    console.log(finalizeDepositData.data);
    if (finalizeDepositData.status === 'successful') {
      onGoToOverviewStep();
    }
  }, [finalizeDepositData, onGoBack, onGoToOverviewStep]);

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
          Check successful transaction details on {NETWORKS[CHIADO_NETWORK_ID].name} for deposit here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      <div className={classes.spinnerWrapper}>
        <Spinner/>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[SPECULAR_NETWORK_ID].blockExplorerUrl}/tx/${finalizeDepositData.data}`}
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
