import { useEffect } from 'react';
import Header from '../shared/header/header.view';
import useTxPendingStyles from './tx-deposit-pending.styles';
import { ReactComponent as LinkIcon } from '../../images/link-icon.svg';
import Spinner from '../shared/spinner/spinner.view';
import { NETWORKS } from '../../chains';

interface TxPendingProps {
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
  onGoBack: () => void;
  onGoToOverviewStep: () => void;
}

function TxPending({ wallet, depositData, onGoBack, onGoToOverviewStep }: TxPendingProps) {
  const classes = useTxPendingStyles();

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

export default TxPending;
