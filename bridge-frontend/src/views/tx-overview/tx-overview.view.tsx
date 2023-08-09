import Header from '../shared/header/header.view';
import useTxOverviewStyles from './tx-overview.styles';
import { ReactComponent as CheckIcon } from '../../images/check-icon.svg';
import LinkIcon from '@mui/icons-material/OpenInNew';
import { NETWORKS } from '../../chains';
import * as React from 'react';
import {
  CHIADO_NETWORK_ID,
  SPECULAR_NETWORK_ID
} from "../../constants";

interface TxOverviewProps {
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
  finalizeTransactionData: {
    status: string;
    data?: string;
  };
  onDisconnectWallet: () => void;
  isDeposit: boolean;
}

function TxOverview({
  wallet,
  transactionData,
  finalizeTransactionData,
  onDisconnectWallet,
  isDeposit,
}: TxOverviewProps) {
  const classes = useTxOverviewStyles();

  var fromNetworkId = SPECULAR_NETWORK_ID;
  var toNetworkId = CHIADO_NETWORK_ID;

  if(isDeposit){
    fromNetworkId = CHIADO_NETWORK_ID;
    toNetworkId = SPECULAR_NETWORK_ID;
  }

  return (
    <div className={classes.txOverview}>
      <Header
        address={wallet.address}
        title={`xDAI → ETH`}
        onDisconnectWallet={onDisconnectWallet}
      />
      <CheckIcon className={classes.checkIcon} />
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[fromNetworkId].blockExplorerUrl}/tx/${transactionData?.data?.hash}`}
          target='_blank'
          rel='noopener noreferrer'
          >
          Check {NETWORKS[fromNetworkId].name}'s transaction details here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
      <div className={classes.buttonGroup}>
        <a
          className={classes.button}
          href={`${NETWORKS[toNetworkId].blockExplorerUrl}/tx/${finalizeTransactionData?.data}`}
          target='_blank'
          rel='noopener noreferrer'
        >
          Check {NETWORKS[toNetworkId].name}'s transaction details here
          <LinkIcon className={classes.buttonIcon} />
        </a>
      </div>
    </div>
  );
}

export default TxOverview;
